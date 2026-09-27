<script setup lang="ts">
import AppHeader from "@/components/app-header.vue";
import UploadFile from "@/components/upload-file.vue";
import ListFiles from "@/components/list-files.vue";
import type { FileInfo } from "@/lib/api";
import { ref } from "vue";

const files = ref<FileInfo[]>([]);

const handleUpload = (file: FileInfo) => files.value.push(file);
</script>

<template>
  <main class="mx-auto max-w-2xl space-y-6 p-6">
    <AppHeader />
    <UploadFile @upload="handleUpload" />
    <ListFiles v-model="files" />
  </main>
</template>
