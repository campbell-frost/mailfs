<script setup lang="ts">
import { listFiles, type FileInfo } from "@/lib/api";
import { onMounted, ref } from "vue";
import { formatBytes } from "@/lib/utils";
import ViewFile from "./view-file.vue";
import { Button } from "@/components/ui/button";

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

const preview = ref<FileInfo | undefined>();

onMounted(fetchFiles);
</script>

<template>
  <ul class="divide-y border-y text-sm">
    <li
      v-for="file in files"
      :key="file.id"
      class="flex items-center gap-3 p-2"
    >
      <div class="flex min-w-0 flex-1 items-center gap-3">
        <span class="flex-1 font-medium">{{ file.filename }}</span>
        <span class="shrink-0 text-muted-foreground">
          {{ formatBytes(file.size) }}
        </span>
      </div>
      <Button size="sm" @click="preview = file">View</Button>
    </li>
  </ul>
  <ViewFile :file="preview" @close="preview = undefined" />
</template>
