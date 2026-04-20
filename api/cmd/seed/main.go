package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azfirazka/gofin-full/api/internal/auth"
	"github.com/azfirazka/gofin-full/api/internal/repository"
)

func main() {
	dsn := flag.String("dsn", "", "Database connection string (required)")
	email := flag.String("email", "", "Admin email (required)")
	password := flag.String("password", "", "Admin password (auto-generated if empty)")
	flag.Parse()

	if *dsn == "" || *email == "" {
		fmt.Fprintln(os.Stderr, "Usage: seed -dsn <DB_DSN> -email <admin@email.com> [-password <secret>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, *dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL connect: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL ping: %v\n", err)
		os.Exit(1)
	}

	// Check if users already exist
	userRepo := repository.NewUserRepository(pool)
	exists, err := userRepo.Exists(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL check existing users: %v\n", err)
		os.Exit(1)
	}
	if exists {
		fmt.Println("SKIP: users already exist in the database.")
		fmt.Println("Use the admin API (POST /api/v1/admin/users) to create additional users.")
		os.Exit(0)
	}

	// Generate password if not provided
	if *password == "" {
		pwd, err := generatePassword(16)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL generate password: %v\n", err)
			os.Exit(1)
		}
		*password = pwd
	}

	hash, err := auth.HashPassword(*password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL hash password: %v\n", err)
		os.Exit(1)
	}

	user, err := userRepo.Create(ctx, *email, hash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL create admin user: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Admin user created ===")
	fmt.Printf("  ID:       %d\n", user.ID)
	fmt.Printf("  Email:    %s\n", user.Email)
	fmt.Printf("  Password: %s\n", *password)
	fmt.Println()
	fmt.Println("IMPORTANT: Save this password. It will not be shown again.")
	fmt.Println("Login via: POST /api/v1/auth/login")
}

func generatePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random password: %w", err)
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}
