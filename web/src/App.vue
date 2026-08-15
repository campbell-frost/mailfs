<script setup lang="ts">
import { ref } from "vue";
import Button from "@/components/ui/button/Button.vue";
import { useDark, useToggle } from "@vueuse/core";
import { Moon, Sun } from "@lucide/vue";
import Input from "@/components/ui/input/Input.vue";

const API_URL = "http://localhost:1738";

const uploadLoading = ref(false);
const uploadError = ref<string | undefined>();

const isDark = useDark();
const toggleDark = useToggle(isDark);

const file = ref<File | null>();
const inputKey = ref(0);

const onFileChange = (e: Event) => {
  file.value = (e.target as HTMLInputElement).files?.[0] ?? null;
};

const upload = async (file: File) => {
  uploadLoading.value = true;
  const form = new FormData();
  form.append("mailfs.file", file);
  try {
    const res = await fetch(`${API_URL}/upload/`, {
      method: "POST",
      body: form,
    });

    if (!res.ok) {
      throw new Error(`bad status: ${res.status}`);
    }
    uploadError.value = undefined;
  } catch (e) {
    uploadError.value = e instanceof Error ? e.message : String(e);
  } finally {
    uploadLoading.value = false;
  }
};
</script>

<template>
  <main class="mx-auto max-w-2xl space-y-6 p-6">
    <div class="flex items-center justify-between">
      <h1 class="font-mono text-lg font-bold">mail.fs</h1>
      <Button
        size="icon-sm"
        variant="ghost"
        :title="isDark ? 'Switch to light' : 'Switch to dark'"
        @click="toggleDark()"
      >
        <Sun v-if="isDark" />
        <Moon v-else />
      </Button>
    </div>
    <div class="flex items-center gap-2">
      <Input
        :key="inputKey"
        type="file"
        class="min-w-0 flex-1 text-sm"
        :disabled="uploadLoading"
        @change="onFileChange"
      />
      <Button
        size="sm"
        :disabled="!file || uploadLoading"
        @click="file && upload(file)"
      >
        {{ uploadLoading ? "uploading..." : "upload" }}
      </Button>
    </div>
    <p v-if="uploadError" class="text-sm text-destructive">
      {{ uploadError }}
    </p>
  </main>
</template>
