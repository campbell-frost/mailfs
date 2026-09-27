import { router } from "@/router";
import { useAuth } from "@/lib/auth";

export const API_URL = "/api";

const apiFetch = async (path: string, init?: RequestInit) => {
  const res = await fetch(`${API_URL}${path}`, init);
  if (res.status === 401) {
    useAuth().user.value = null;
    router.push({ name: "login" });
  }
  return res;
};

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
  const res = await apiFetch("/list/");
  if (!res.ok) {
    throw new Error(`failed to load files: ${res.status}`);
  }
  return await res.json();
};

export const uploadFile = async (file: File): Promise<FileInfo> => {
  const form = new FormData();
  form.append("mailfs.file", file);

  const res = await apiFetch("/upload/", {
    method: "POST",
    body: form,
  });

  if (!res.ok) {
    throw new Error(`bad status: ${res.status}`);
  }
  return await res.json();
};

export const downloadFile = async (id: string): Promise<Response> => {
  const res = await apiFetch(`/download/${id}`);
  if (!res.ok) {
    throw new Error(`bad status: ${res.status}`);
  }
  return res;
};

export const deleteFile = async (id: string): Promise<void> => {
  const res = await apiFetch(`/delete/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    throw new Error(`bad status: ${res.status} ${res.statusText}`);
  }
};

export type FileStatus = {
  status: "pending" | "processing" | "stored" | "failed";
  processedChunks: number;
};

export const getFileStatus = async (id: string): Promise<FileStatus> => {
  const res = await apiFetch(`/status/${id}`);
  if (!res.ok) {
    throw new Error(`bad status: ${res.status} ${res.statusText}`);
  }

  const json = await res.json();

  // TODO: use zod to validate response
  return json as FileStatus;
};
