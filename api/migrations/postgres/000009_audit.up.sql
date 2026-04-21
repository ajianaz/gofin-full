-- Audit logs for tracking financial mutations
CREATE TABLE IF NOT EXISTS audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         UUID NOT NULL,
    user_group_id   UUID NOT NULL,
    action          VARCHAR(255) NOT NULL,
    entity_type     VARCHAR(255) NOT NULL DEFAULT '',
    entity_id       UUID NOT NULL,
    old_value       TEXT,
    new_value       TEXT,
    ip_address      VARCHAR(45),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_group ON audit_logs(user_group_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);
