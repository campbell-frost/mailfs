<script setup lang="ts">
import { useDark, useToggle } from "@vueuse/core";
import { CircleUser, LogOut, Moon, Sun } from "@lucide/vue";
import { useRouter } from "vue-router";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu/index.ts";
import { useAuth } from "@/lib/auth";

const isDark = useDark();
const toggleDark = useToggle(isDark);

const router = useRouter();
const { user, logout } = useAuth();

const handleLogout = async () => {
  await logout();
  router.push({ name: "login" });
};
</script>

<template>
  <div class="flex items-center justify-between">
    <h1 class="font-mono text-lg font-bold">mail.fs</h1>
    <div class="flex items-center gap-1">
      <Button
        size="icon-sm"
        variant="ghost"
        :title="isDark ? 'Switch to light' : 'Switch to dark'"
        @click="toggleDark()"
      >
        <Sun v-if="isDark" />
        <Moon v-else />
      </Button>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button size="icon-sm" variant="ghost">
            <CircleUser />
          </Button>
        </DropdownMenuTrigger>

        <DropdownMenuContent align="end" class="w-48">
          <DropdownMenuLabel class="flex flex-col">
            <span>{{ user?.name }}</span>
            <span class="truncate text-xs font-normal text-muted-foreground">
              {{ user?.email }}
            </span>
          </DropdownMenuLabel>

          <DropdownMenuSeparator />

          <DropdownMenuItem @select="handleLogout">
            <LogOut class="mr-2 size-4" />
            Sign out
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </div>
</template>
