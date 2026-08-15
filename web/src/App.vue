<script setup lang="ts">
import { ref } from "vue";
import Button from "@/components/ui/button/Button.vue";

const API_URL = "http://localhost:1738";

const uploadLoading = ref(false);
const uploadError = ref<string | undefined>();
const uploadResult = ref<{ id: string }>();

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
  <div class="min-h-screen flex items-center justify-center">
    <div class="flex flex-col">
      <Button @click="upload">upload</Button>
      <div>
        <p v-if="uploadLoading">Loading...</p>
        <p v-if="uploadError">Error {{ uploadError }}</p>
        <p v-if="uploadResult">Result {{ uploadResult }}</p>
      </div>
    </div>
  </div>
</template>
