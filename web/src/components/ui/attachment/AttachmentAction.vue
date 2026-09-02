<script setup lang="ts">
import type { PrimitiveProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import type { ButtonVariants } from "../button"
import Button from "../button/Button.vue"

interface Props extends PrimitiveProps {
  class?: HTMLAttributes["class"]
  variant?: ButtonVariants["variant"]
  size?: ButtonVariants["size"]
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: "ghost",
  size: "icon-xs",
  disabled: false,
})
</script>

<template>
  <Button
    data-slot="attachment-action"
    :variant="variant"
    :size="size"
    :disabled
    >
    <slot />
  </Button>
</template>
