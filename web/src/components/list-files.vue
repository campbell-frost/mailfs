<script setup lang="ts">
import { listFiles, type FileInfo } from "@/lib/api";
import { onMounted, ref } from "vue";
import { formatBytes } from "@/lib/utils";
import ViewFile from "@/components/view-file.vue";
import { Button } from "@/components/ui/button";
import DeleteFile from "@/components/delete-file.vue";

const files = defineModel<FileInfo[]>({
  default: () => [],
});

const isLoading = ref(false);
const error = ref<Error | null>(null);

const fetchFiles = async () => {
  isLoading.value = true;
  error.value = null;
  try {
    files.value = await listFiles();
  } catch (e) {
    error.value = e instanceof Error ? e : new Error("An error occurred");
  } finally {
    isLoading.value = false;
  }
};

const preview = ref<FileInfo | undefined>();
const fileToDelete = ref<FileInfo | undefined>();

const handleDelete = (id: string | undefined) => {
  fileToDelete.value = undefined;
  // id is undefined if the delete dialog was cancelled, dont remove the file from the UI.
  if (id == null) return;
  files.value = files.value.filter(f => f.id !== id);
};

onMounted(fetchFiles);
</script>

<template>
  <p v-if="isLoading">Loading...</p>
  <p v-if="error" class="text-sm text-destructive">
    {{ error }}
  </p>
  <div
    v-if="files.length === 0"
    class="flex items-center justify-center border p-4 rounded-md"
  >
    <p>No files have been uploaded yet.</p>
  </div>
  <ul v-else class="divide-y border-y text-sm">
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
      <Button size="sm" @click="fileToDelete = file">Delete</Button>
    </li>
  </ul>
  <ViewFile :file="preview" @close="preview = undefined" />
  <DeleteFile :file="fileToDelete" @close="handleDelete" />
</template>
