<script setup lang="ts">
import { computed, ref } from "vue";
import { getFileStatus, uploadFile, type FileInfo } from "@/lib/api";
import { formatBytes } from "@/lib/utils";
import {
  Attachment,
  AttachmentContent,
  AttachmentMedia,
  AttachmentDescription,
  AttachmentTitle,
  AttachmentAction,
  AttachmentActions,
} from "@/components/ui/attachment/index.ts";
import {
  Check,
  FileTextIcon,
  Loader2,
  Upload,
  XIcon,
  type LucideIcon,
} from "@lucide/vue";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip/index.ts";

type UploadState = "idle" | "uploading" | "processing" | "error" | "done";

const emit = defineEmits<{
  upload: [file: FileInfo];
}>();

const file = ref<File | null>(null);
const fileInput = ref<HTMLInputElement | null>(null);
const fileState = ref<UploadState>("idle");
const uploadError = ref<string>();

const totalChunks = ref(0);
const proccessedChunks = ref(0);

const selectFile = () => fileInput.value?.click();
const removeFile = () => {
  fileInput.value = null;
  file.value = null;
};

const onFileChange = (e: Event) => {
  const input = e.target as HTMLInputElement;

  file.value = input.files?.[0] ?? null;
  fileState.value = "idle";
  uploadError.value = undefined;
};

const upload = async () => {
  if (!file.value) return;

  fileState.value = "uploading";
  uploadError.value = undefined;

  try {
    const fi = await uploadFile(file.value);
    totalChunks.value = fi.chunks;
    await pollStatus(fi.id, 2000);
    emit("upload", fi);
  } catch (e) {
    fileState.value = "error";
    uploadError.value = e instanceof Error ? e.message : String(e);
  }
};

const pollStatus = async (id: string, ms: number) => {
  while (true) {
    try {
      const fs = await getFileStatus(id);
      switch (fs.status) {
        case "pending":
          fileState.value = "uploading";
          break;

        case "processing":
          fileState.value = "processing";
          proccessedChunks.value = fs.processedChunks;
          break;

        case "stored":
          proccessedChunks.value = 0;
          fileState.value = "done";
          file.value = null;
          return;

        case "failed":
          fileState.value = "error";
          uploadError.value = "File processing failed";
          return;

        default:
          fileState.value = "error";
          proccessedChunks.value = 0;
          uploadError.value = `Unknown file status: ${fs.status}`;
          return;
      }

      await new Promise(r => setTimeout(r, ms));
    } catch (e) {
      fileState.value = "error";
      uploadError.value = e instanceof Error ? e.message : String(e);
    }
  }
};

const fileIcon = computed<LucideIcon>(() => {
  switch (fileState.value) {
    case "processing":
      return Loader2;

    case "error":
      return XIcon;

    case "done":
      return Check;

    default:
      return FileTextIcon;
  }
});

const fileIconClass = computed(() =>
  fileState.value === "processing" ? "animate-spin" : "",
);
</script>

<template>
  <div
    class="cursor-pointer rounded-lg border border-dashed p-6 text-center hover:bg-muted/50"
    @click="selectFile"
  >
    <p class="text-muted-foreground">Click to select a file</p>
    <input ref="fileInput" type="file" class="hidden" @change="onFileChange" />
  </div>

  <Attachment v-if="file" class="flex w-full" :state="fileState">
    <AttachmentMedia>
      <component :is="fileIcon" :class="fileIconClass" />
    </AttachmentMedia>

    <AttachmentContent>
      <AttachmentTitle>{{ file.name }}</AttachmentTitle>
      <AttachmentDescription>
        <span>{{ formatBytes(file.size) }}</span>
        <span v-if="fileState !== 'idle'">
          ({{ proccessedChunks }} / {{ totalChunks }}) chunks processed
        </span>
      </AttachmentDescription>
    </AttachmentContent>

    <AttachmentActions
      class="flex items-center gap-1 rounded-md border bg-muted/30 p-1"
    >
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger as-child>
            <AttachmentAction
              :disabled="fileState === 'uploading'"
              @click="upload"
            >
              <Upload />
            </AttachmentAction>
          </TooltipTrigger>

          <TooltipContent> Upload File </TooltipContent>
        </Tooltip>
        <div class="h-3 w-px bg-border" />

        <Tooltip>
          <TooltipTrigger as-child>
            <AttachmentAction
              :disabled="fileState === 'uploading'"
              @click="removeFile"
            >
              <XIcon />
            </AttachmentAction>
          </TooltipTrigger>

          <TooltipContent> Clear </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </AttachmentActions>
  </Attachment>

  <p v-if="uploadError" class="text-sm text-destructive">
    {{ uploadError }}
  </p>
</template>
