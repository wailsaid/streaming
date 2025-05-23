<script setup>
import {  Search, Bell, User, Video,  History, Upload, LogOut, Library, Home, Compass, Menu } from "lucide-vue-next"
import Button from "./ui/button/Button.vue";
import { RouterLink } from "vue-router";
import Input from "@/components/ui/input/Input.vue"

import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import SheetTrigger from "./ui/sheet/SheetTrigger.vue";
import Sheet from "./ui/sheet/Sheet.vue";
import SheetContent from "./ui/sheet/SheetContent.vue";

// Add mock notifications data
const mockNotifications = [
    {
        id: 1,
        title: "New video from CodeWithMosh",
        message: "Advanced TypeScript Patterns just dropped!",
        time: "2 hours ago",
        read: false,
    },
    {
        id: 2,
        title: "Fireship replied to your comment",
        message: "Thanks for the feedback!",
        time: "1 day ago",
        read: false,
    },
    {
        id: 3,
        title: "Your video is done processing",
        message: "Your upload 'React Hooks Tutorial' is now available",
        time: "3 days ago",
        read: true,
    },
]

/*
export default function Navbar() {
const [searchQuery, setSearchQuery] = useState("")
const router = useRouter()
const pathname = usePathname()
const isHomePage = pathname === "/"

const handleSearch = (e: React.FormEvent) => {
e.preventDefault()
if (searchQuery.trim()) {
  // In a real app, this would navigate to search results
  // For now, we'll just log the query
  console.log("Searching for:", searchQuery)
  alert(`Search functionality would search for: ${searchQuery}`)
}
} */
</script>
<template>
    <header
        class="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div class="container flex h-14 items-center max-w-7xl mx-auto px-4">
            <div class="flex items-center mr-2 md:mr-6">
                 <Sheet>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" class="md:hidden">
                <Menu class="h-5 w-5" />
                <span class="sr-only">Toggle menu</span>
              </Button>
            </SheetTrigger>
            <SheetContent side="left" class="w-[300px] sm:w-[350px]">
              <nav class="flex flex-col gap-4 mt-8">
                <RouterLink to="/" class="flex items-center gap-2 py-2 font-medium">
                  <Home class="h-5 w-5" />
                  Home
                </RouterLink>
                <RouterLink to="/trending" class="flex items-center gap-2 py-2 font-medium">
                  <Compass class="h-5 w-5" />
                  Trending
                </RouterLink>
                <RouterLink to="/subscriptions" class="flex items-center gap-2 py-2 font-medium">
                  <Library class="h-5 w-5" />
                  Subscriptions
                </RouterLink>
                <RouterLink to="/history" class="flex items-center gap-2 py-2 font-medium">
                  <History class="h-5 w-5" />
                  History
                </RouterLink>
              </nav>
            </SheetContent>
          </Sheet>


                <RouterLink to="/" class="flex items-center space-x-2">
                    <Video class="h-6 w-6 text-red-600" />
                    <span class="font-bold text-lg hidden sm:inline-block">YouClone</span>
                </RouterLink>
            </div>

            <div class="flex-1 flex justify-center px-2">
                <form onSubmit="" class="w-full max-w-lg flex">
                    <Input type="search" placeholder="Search" class="rounded-r-none" value="" onChange="" />
                    <Button type="submit" variant="secondary" class="rounded-l-none">
                        <Search class="h-4 w-4" />
                        <span class="sr-only">Search</span>
                    </Button>
                </form>
            </div>

            <div class="flex items-center gap-2">
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon" class="hidden md:flex relative">
                            <Bell class="h-5 w-5" />
                            <span v-if="mockNotifications.some((n) => !n.read)"
                                class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
                            <span class="sr-only">Notifications</span>
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" class="w-80">
                        <DropdownMenuLabel>Notifications</DropdownMenuLabel>
                        <DropdownMenuSeparator />

                        <template v-if="mockNotifications.length > 0">

                            <DropdownMenuItem v-for="notification in mockNotifications" key="" class="">
                                <div class="flex flex-col gap-1">
                                    <span class="font-medium">{{ notification.title }}</span>
                                    <span class="text-sm text-muted-foreground">{{ notification.message }}</span>
                                    <span class="text-xs text-muted-foreground">{{ notification.time }}</span>
                                </div>
                            </DropdownMenuItem>

                            <DropdownMenuSeparator />
                            <DropdownMenuItem class="cursor-pointer justify-center font-medium">
                                See all notifications
                            </DropdownMenuItem>

                        </template>

                        <div v-else class="p-4 text-center">
                            <p class="text-muted-foreground">No notifications</p>
                        </div>
                    </DropdownMenuContent>
                </DropdownMenu>

                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon">
                            <User class="h-5 w-5" />
                            <span class="sr-only">Profile</span>
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" class="w-56">
                        <DropdownMenuLabel>My Account</DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem asChild>
                            <RouterLink to="/profile" class="cursor-pointer w-full">
                                <User class="mr-2 h-4 w-4" />
                                <span>Profile</span>
                            </RouterLink>
                        </DropdownMenuItem>
                        <DropdownMenuItem asChild>
                            <RouterLink to="/profile/upload" class="cursor-pointer w-full">
                                <Upload class="mr-2 h-4 w-4" />
                                <span>Upload Video</span>
                            </RouterLink>
                        </DropdownMenuItem>
                        <DropdownMenuItem asChild>
                            <RouterLink to="/profile/history" class="cursor-pointer w-full">
                                <History class="mr-2 h-4 w-4" />
                                <span>Watch History</span>
                            </RouterLink>
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem>
                            <LogOut class="mr-2 h-4 w-4" />
                            <span>Log out</span>
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </div>

        <nav v-if="false && 'isHomePage'" class="hidden md:flex overflow-auto whitespace-nowrap border-t px-4 py-2">
            <div class="flex gap-2 mx-auto max-w-7xl">
                <RouterLink to="/">
                    <Button variant="secondary" size="sm" class="rounded-full">
                        All
                    </Button>
                </RouterLink>
                <RouterLink to="/?category=music">
                    <Button variant="ghost" size="sm" class="rounded-full">
                        Music
                    </Button>
                </RouterLink>
                <RouterLink to="/?category=gaming">
                    <Button variant="ghost" size="sm" class="rounded-full">
                        Gaming
                    </Button>
                </RouterLink>
                <RouterLink to="/?category=news">
                    <Button variant="ghost" size="sm" class="rounded-full">
                        News
                    </Button>
                </RouterLink>
                <RouterLink to="/?category=movies">
                    <Button variant="ghost" size="sm" class="rounded-full">
                        Movies
                    </Button>
                </RouterLink>
                <RouterLink to="/?category=programming">
                    <Button variant="ghost" size="sm" class="rounded-full">
                        Programming
                    </Button>
                </RouterLink>
                <RouterLink to="/?category=design">
                    <Button variant="ghost" size="sm" class="rounded-full">
                        Design
                    </Button>
                </RouterLink>
                <RouterLink to="/?category=education">
                    <Button variant="ghost" size="sm" class="rounded-full">
                        Education
                    </Button>
                </RouterLink>
            </div>
        </nav>

    </header>
</template>