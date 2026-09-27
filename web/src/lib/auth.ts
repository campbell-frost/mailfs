import { ref } from "vue";

export type User = {
  id: number;
  email: string;
  name: string;
};

const user = ref<User | null>(null);
let loaded: Promise<void> | null = null;

const fetchSession = async () => {
  const res = await fetch("/api/session");
  user.value = res.ok ? await res.json() : null;
};

const logout = async () => {
  await fetch("/api/session", { method: "DELETE" });
  user.value = null;
};

export const useAuth = () => {
  loaded ??= fetchSession();
  return { user, loaded, logout };
};
