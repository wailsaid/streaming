/* 
import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
import type { Ref } from 'vue'
*/
export function cn(...inputs: any[]): string {
  return inputs
    .flat(Infinity) // Flatten nested arrays
    .filter(Boolean) // Remove falsy values like null, undefined, false, 0, ''
    .join(' ')
} 