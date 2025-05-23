<template>
  <component
    :is="asChild ? 'slot' : 'button'"
    :class="computedClass"
    v-bind="attrs"
    ref="btnRef"
  >
    <slot />
  </component>
</template>

<script setup lang="ts">
import { computed, defineProps, ref, useAttrs } from 'vue'
import { cn } from '../lib/utils'

const props = defineProps<{
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  asChild?: boolean
  class?: string
}>()

const attrs = useAttrs()
const btnRef = ref<HTMLElement | null>(null)

const variantClasses: Record<string, string> = {
  default: 'bg-primary text-primary-foreground hover:bg-primary/90',
  destructive: 'bg-destructive text-destructive-foreground hover:bg-destructive/90',
  outline: 'border border-input bg-background hover:bg-accent hover:text-accent-foreground',
  secondary: 'bg-secondary text-secondary-foreground hover:bg-secondary/80',
  ghost: 'hover:bg-accent hover:text-accent-foreground',
  link: 'text-primary underline-offset-4 hover:underline',
}

const sizeClasses: Record<string, string> = {
  default: 'h-10 px-4 py-2',
  sm: 'h-9 rounded-md px-3',
  lg: 'h-11 rounded-md px-8',
  icon: 'h-10 w-10',
}

const baseClass = `inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0`

const computedClass = computed(() =>
  cn(
    baseClass,
    variantClasses[props.variant || 'default'],
    sizeClasses[props.size || 'default'],
    props.class
  )
)
</script>
