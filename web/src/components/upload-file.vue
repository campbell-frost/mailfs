<script setup lang="ts">
import { uploadFile } from "@/lib/api";
import { ref } from "vue";
import { Input } from "./ui/input";
import { Button } from "./ui/button";

const uploadLoading = ref(false);
const uploadError = ref<string | undefined>();

const file = ref<File | null>();
const inputKey = ref(0);

const onFileChange = (e: Event) => {
  file.value = (e.target as HTMLInputElement).files?.[0] ?? null;
};

const upload = async () => {
  if (!file.value) return;
  uploadLoading.value = true;
  try {
    await uploadFile(file.value);
    uploadError.value = undefined;
  } catch (e) {
    uploadError.value = e instanceof Error ? e.message : String(e);
  } finally {
    uploadLoading.value = false;
  }
};
</script>
<template>
  <div class="flex items-center gap-2">
    <Input
      :key="inputKey"
      type="file"
      class="min-w-0 flex-1 text-sm"
      :disabled="uploadLoading"
      @change="onFileChange"
    />
    <Button :disabled="!file || uploadLoading" @click="file && upload()">
      {{ uploadLoading ? "uploading..." : "upload" }}
    </Button>
  </div>
  <p v-if="uploadError" class="text-sm text-destructive">
    {{ uploadError }}
  </p>
</template>
