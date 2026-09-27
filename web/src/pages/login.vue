<script setup lang="ts">
import { Button } from "@/components/ui/button";
</script>

<template>
  <div class="flex min-h-svh items-center justify-center p-6">
    <div class="flex flex-col items-center gap-4">
      <Button as-child>
        <a href="/auth/google/login">Sign in with Google</a>
      </Button>
    </div>
  </div>
</template>
