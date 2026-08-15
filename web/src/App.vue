<script setup lang="ts">
import { ref } from "vue";
import Button from "@/components/ui/button/Button.vue";
import { useDark, useToggle } from "@vueuse/core";
import { Moon, Sun } from "@lucide/vue";

const API_URL = "http://localhost:1738";

const uploadLoading = ref(false);
const uploadError = ref<string | undefined>();
const uploadResult = ref<{ id: string }>();

const isDark = useDark();
const toggleDark = useToggle(isDark);

const upload = async () => {
  uploadLoading.value = true;
  try {
    const res = await fetch(`${API_URL}/upload/`, {
      method: "POST",
      body: JSON.stringify({
        id: "hi",
      }),
    });

    if (!res.ok) {
      throw new Error(`bad status: ${res.status}`);
    }

    uploadResult.value = await res.json();
    uploadError.value = undefined;
  } catch (e) {
    uploadError.value = e instanceof Error ? e.message : String(e);
  } finally {
    uploadLoading.value = false;
  }
};
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
    <div class="flex flex-col">
      <Button @click="upload">upload</Button>
      <div>
        <p v-if="uploadLoading">Loading...</p>
        <p v-if="uploadError">Error {{ uploadError }}</p>
        <p v-if="uploadResult">Result {{ uploadResult }}</p>
      </div>
    </div>
  </main>
</template>
