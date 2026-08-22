<script setup lang="ts">
import { useDark, useToggle } from "@vueuse/core";
import { Moon, Sun } from "@lucide/vue";

const isDark = useDark();
const toggleDark = useToggle(isDark);
</script>

<template>
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
</template>
