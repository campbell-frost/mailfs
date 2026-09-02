export const API_URL = "http://localhost:1738";

export type Status = "pending" | "processing" | "stored" | "failed";

export type FileInfo = {
  id: string;
  filename: string;
  size: number;
  status: Status;
  createdAt: string;
  chunks: number;
};

export const listFiles = async (): Promise<FileInfo[]> => {
  const res = await fetch(`${API_URL}/list/`);
  if (!res.ok) {
    throw new Error(`failed to load files: ${res.status}`);
  }
  return await res.json();
};

export const uploadFile = async (file: File): Promise<FileInfo> => {
  const form = new FormData();
  form.append("mailfs.file", file);

  const res = await fetch(`${API_URL}/upload/`, {
    method: "POST",
    body: form,
  });

  if (!res.ok) {
    throw new Error(`bad status: ${res.status}`);
  }
  return await res.json();
};

export const downloadFile = async (id: string): Promise<Response> => {
  const res = await fetch(`${API_URL}/download/${id}`);
  if (!res.ok) {
    throw new Error(`bad status: ${res.status}`);
  }
  return res;
};

export const deleteFile = async (id: string): Promise<void> => {
  const res = await fetch(`${API_URL}/delete/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    throw new Error(`bad status: ${res.status} ${res.statusText}`);
  }
};

export type FileStatus = "pending" | "processing" | "stored" | "failed";

export const getStatus = async (id: string): Promise<FileStatus> => {
  const res = await fetch(`${API_URL}/status/${id}`);
  if (!res.ok) {
    throw new Error(`bad status: ${res.status} ${res.statusText}`);
  }

  const json = await res.json();

  // TODO: use zod to validate response
  return json.status as FileStatus;
};
