<script setup lang="ts">
import { deleteFile, type FileInfo } from "@/lib/api";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
} from "@/components/ui/dialog";
import { ref } from "vue";
import { Button } from "./ui/button";

interface Props {
  file: FileInfo | undefined;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  close: [id: string | undefined];
}>();

const loading = ref(false);
const error = ref("");

const confirmDelete = async () => {
  loading.value = true;
  try {
    if (!props.file) return;
    const id = props.file.id;
    await deleteFile(id);
    emit("close", id);
  } catch (e) {
    error.value = "Failed to delete file";
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <Dialog
    :open="!!file"
    @update:open="open => !open && emit('close', undefined)"
  >
    <DialogContent v-if="file">
      <DialogTitle>Delete {{ file.filename }}</DialogTitle>
      <DialogDescription>
        This throws away the {{ file.chunks }}
        {{ file.chunks === 1 ? "email" : "emails" }} holding the data. This
        cannot be undone.
      </DialogDescription>
      <p v-if="error" class="text-destructive text-sm">{{ error }}</p>
      <DialogFooter>
        <Button size="sm" @click="emit('close', undefined)">Cancel</Button>
        <Button
          size="sm"
          variant="destructive"
          :disabled="loading"
          @click="confirmDelete"
        >
          {{ loading ? "Deleting..." : "Delete" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
