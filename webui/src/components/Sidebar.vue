<script setup lang="ts">
import {
    Home,
    Library,
    History,
    Clock,
    ThumbsUp,
    Film,
    Flame,
    Music,
    Gamepad2,
    Newspaper,
    Lightbulb,
} from "lucide-vue-next"

import { RouterLink, useRoute } from "vue-router";
import ScrollArea from "./ui/scroll-area/ScrollArea.vue"
import { cn } from "../lib/utils";
import Button from "./ui/button/Button.vue";

const useroute = useRoute();

const mainLinks = [
    { href: "/", label: "Home", icon: Home },
    { href: "/trending", label: "Trending", icon: Flame },
    { href: "/subscriptions", label: "Subscriptions", icon: Library },
]

const libraryLinks = [
    { href: "/profile/history", label: "History", icon: History },
    { href: "/profile/watch-later", label: "Watch Later", icon: Clock },
    { href: "/profile/liked-videos", label: "Liked Videos", icon: ThumbsUp },
]

const exploreLinks = [
    { href: "/?category=movies", label: "Movies", icon: Film },
    { href: "/?category=music", label: "Music", icon: Music },
    { href: "/?category=gaming", label: "Gaming", icon: Gamepad2 },
    { href: "/?category=news", label: "News", icon: Newspaper },
    { href: "/?category=learning", label: "Learning", icon: Lightbulb },
]

</script>
<template>
        <aside class="hidden md:block w-60 fixed top-14 left-0 bottom-0 border-r bg-background z-30">
            <ScrollArea className="h-full py-4">
                <div class="px-3 space-y-4">
                    
                  
                        <div class="space-y-1">
                            <RouterLink v-for="link in mainLinks" :key="link.label" :to="link.href">
                                <Button variant="ghost"
                                    :class="cn('w-full !justify-start', useroute.path === link.href && 'bg-muted')">
                                    <component class="mr-2 h-4 w-4" :is="link.icon" />
                                    <span class="">{{ link.label }}</span>
                                </Button>
                            </RouterLink>

                        </div>
                   

                    <hr />

                    <div class="space-y-1">
                        <h3 class="px-4 text-sm font-medium">Library</h3>

                        <RouterLink v-for="link in libraryLinks" :key="link.label" :to="link.href">
                            <Button variant="ghost"
                                :class="cn('w-full !justify-start', useroute.path === link.href && 'bg-muted')">
                                <component class="mr-2 h-4 w-4" :is="link.icon" />
                                {{ link.label }}
                            </Button>
                        </RouterLink>

                    </div>

                    <hr />

                    <div class="space-y-1">
                        <h3 class="px-4 text-sm font-medium">Explore</h3>

                        <RouterLink v-for="link in exploreLinks" :key="link.label" :to="link.href">
                            <Button variant="ghost"
                                :class="cn('w-full !justify-start', useroute.path === link.href && 'bg-muted')">
                                <component class="mr-2 h-4 w-4" :is="link.icon" />
                                {{ link.label }}
                            </Button>
                        </RouterLink>

                    </div>

                    <hr/>
                    
                    <div class="px-4 text-xs text-muted-foreground">
                        <p class="mb-2">© 2023 YouClone</p>
                        <div class="flex flex-wrap gap-x-2">
                            <RouterLink to="#" class="hover:underline">
                                Terms
                            </RouterLink>
                            <RouterLink to="#" class="hover:underline">
                                Privacy
                            </RouterLink>
                            <RouterLink to="#" class="hover:underline">
                                Policy & Safety
                            </RouterLink>
                        </div>
                    </div>
                </div>
            </ScrollArea>
        </aside>

</template>
<style scoped></style>