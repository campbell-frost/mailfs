<script setup lang="ts">
import Button from "@/components/ui/button/Button.vue";
import { useDark, useToggle } from "@vueuse/core";
import { Moon, Sun } from "@lucide/vue";
import UploadFile from "./components/upload-file.vue";
import ListFiles from "./components/list-files.vue";

const isDark = useDark();
const toggleDark = useToggle(isDark);
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
    <UploadFile />
    <ListFiles />
  </main>
</template>
