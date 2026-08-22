<script setup lang="ts">
import type { FileInfo } from "@/lib/api";
import { Dialog, DialogContent, DialogHeader } from "@/components/ui/dialog";

interface Props {
  file: FileInfo | undefined;
}

const props = defineProps<Props>();
const emit = defineEmits<{ close: [] }>();
</script>

<template>
  <Dialog :open="!!file" @update:open="open => !open && emit('close')">
    <DialogContent v-if="file">
      <DialogHeader>{{ file.filename }}</DialogHeader>
    </DialogContent>
  </Dialog>
</template>
