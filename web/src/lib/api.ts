export const API_URL = "http://localhost:1738";

export type Status = "pending" | "processing" | "stored" | "failed";

export type FileInfo = {
  id: string;
  filename: string;
  size: number;
  status: Status;
  createdAt: string;
  chunks: number;
}

export const listFiles = async (): Promise<FileInfo[]> => {
  const res = await fetch(`${API_URL}/list/`);
  if (!res.ok) {
    throw new Error(`failed to load files: ${res.status}`);
  }
  return await res.json();
};

export const uploadFile = async (file: File) => {
  const form = new FormData();
  form.append("mailfs.file", file);

  const res = await fetch(`${API_URL}/upload/`, {
    method: "POST",
    body: form,
  });

  if (!res.ok) {
    throw new Error(`bad status: ${res.status}`);
  }
};
