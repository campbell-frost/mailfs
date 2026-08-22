<script setup lang="ts">
import { listFiles, type FileInfo } from "@/lib/api";
import { onMounted, ref } from "vue";
import { formatBytes } from "@/lib/utils";

const files = ref<FileInfo[]>([]);
const isLoading = ref(false);
const error = ref<Error | null>(null);

const fetchFiles = async () => {
  try {
    isLoading.value = true;
    error.value = null;
    files.value = await listFiles();
  } catch (e) {
    error.value = e instanceof Error ? e : new Error("An error occurred");
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchFiles);
</script>

<template>
  <ul class="divide-y border-y text-sm">
    <li v-for="file in files" :key="file.id" class="flex p-2">
      <span class="flex-1 truncate">{{ file.filename }}</span>
      <span>{{ formatBytes(file.size) }}</span>
    </li>
  </ul>
</template>
