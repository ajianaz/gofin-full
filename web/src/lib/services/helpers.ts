export interface JsonApiItem {
	id: string;
	attributes: Record<string, unknown>;
}

export interface JsonApiSingle {
	data: JsonApiItem;
}

export interface JsonApiMany {
	data: JsonApiItem[] | null;
}

export function unwrapOne<T>(res: JsonApiSingle): T & { id: string } {
	return { id: res.data.id, ...res.data.attributes } as T & { id: string };
}

export function unwrapMany<T>(res: JsonApiMany): (T & { id: string })[] {
	return (res.data ?? []).map((item) => ({ id: item.id, ...item.attributes } as T & { id: string }));
}
