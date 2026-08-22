<script setup lang="ts">
import { downloadFile, type FileInfo } from "@/lib/api";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { computed, ref, watch } from "vue";

interface Props {
  file: FileInfo | undefined;
}

const props = defineProps<Props>();
const emit = defineEmits<{ close: [] }>();

const error = ref("");
const loading = ref(false);
const url = ref("");
const text = ref("");

const name = computed(() => props.file?.filename ?? "");
const isImage = computed(() => /\.(png|jpe?g|gif|webp)$/i.test(name.value));
const isPdf = computed(() => /\.pdf$/i.test(name.value));

watch(
  () => props.file,
  async f => {
    loading.value = true;
    error.value = "";
    try {
      if (url.value) URL.revokeObjectURL(url.value);
      text.value = "";
      url.value = "";
      if (!f) return;

      const r = await downloadFile(f.id);

      if (isImage.value) {
        url.value = URL.createObjectURL(await r.blob());
      } else if (isPdf.value) {
        const pdf = new Blob([await r.blob()], { type: "application/pdf" });
        url.value = URL.createObjectURL(pdf);
      } else {
        text.value = (await r.text()).slice(0, 173800);
      }
    } catch (e) {
      error.value = `Could not download file: ${f?.filename} ${e}`;
    } finally {
      loading.value = false;
    }
  },
);
</script>

<template>
  <Dialog :open="!!file" @update:open="open => !open && emit('close')">
    <DialogContent v-if="file" class="sm:max-w-3xl">
      <DialogHeader class="min-w-0">
        <DialogTitle class="pr-6 wrap-break-word">
          {{ file.filename }}
        </DialogTitle>
      </DialogHeader>
      <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
      <div v-if="loading" class="flex justify-center items-center h-full">
        <p class="text-sm text-muted">Loading...</p>
      </div>

      <img
        v-if="isImage"
        :src="url"
        :alt="file.filename"
        class="max-h-[70vh] w-full"
      />

      <iframe
        v-else-if="isPdf"
        :src="url"
        :alt="file.filename"
        class="h-[70vh] w-full"
      />
      <pre
        v-else
        class="max-h-[70vh] overflow-auto text-xs whitespace-pre-wrap"
      >
        {{ text }}
        </pre>
    </DialogContent>
  </Dialog>
</template>
