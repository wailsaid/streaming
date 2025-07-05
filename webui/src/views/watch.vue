<script setup lang="ts">
import Button from '@/components/ui/button/Button.vue';
import VideoPlayer from '@/components/videoPlayer.vue';
import { ThumbsDown, ThumbsUp } from 'lucide-vue-next';
import { RouterLink, useRoute } from 'vue-router';

const route = useRoute()
</script>
<template>

    
<div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mt-4">
      <div class="lg:col-span-2 space-y-4">
        <VideoPlayer />

     <!--    <div>
          <h1 class="text-xl font-bold mt-2">{video.title}</h1>

          <div class="flex flex-col sm:flex-row sm:items-center justify-between mt-4 gap-4">
            <div class="flex items-center gap-3">
              <RouterLink to="/channel/${video.channelId}">
                <img
                  src=""
                  alt="video.channelName"
                  width="48"
                  height="48"
                  class="rounded-full"
                />
              </RouterLink>
              <div>
                <RouterLink to="/channel/${video.channelId}" class="font-medium hover:text-blue-600">
                  {video.channelName}
                </RouterLink>
                <p class="text-sm text-muted-foreground">1.2M subscribers</p>
              </div>
              <Button class="ml-4">Subscribe</Button>
            </div>

            <div class="flex items-center gap-2">
              <div class="flex rounded-full bg-muted">
                <Button variant="ghost" size="sm" class="rounded-l-full rounded-r-none">
                  <ThumbsUp class="h-4 w-4 mr-2" />
                  12K
                </Button>
                <Button variant="ghost" size="sm" class="rounded-l-none rounded-r-full">
                  <ThumbsDown class="h-4 w-4" />
                </Button>
              </div>

              <Button variant="ghost" size="sm" class="rounded-full">
                <Share class="h-4 w-4 mr-2" />
                Share
              </Button>

              <Button variant="ghost" size="sm" class="rounded-full">
                <Save class="h-4 w-4 mr-2" />
                Save
              </Button>

              <Button variant="ghost" size="icon" class="rounded-full">
                <MoreHorizontal class="h-4 w-4" />
              </Button>
            </div>
          </div>

          <div class="mt-4 bg-muted p-3 rounded-lg">
            <div class="flex items-center gap-2 text-sm">
              <span>{formatViewCount(video.views)} views</span>
              <span>•</span>
              <span>{video.uploadedAt}</span>
            </div>
            <p class="mt-2 whitespace-pre-line">{video.description}</p>
          </div>
        </div>

        <CommentSection comments="comments" />
      </div>

      <div class="space-y-4">
        <h3 class="font-medium">Related videos</h3>
        <div class="space-y-3">
          {relatedVideos.map((video) => (
            <VideoCard key="video.id} video="video} horizontal />
          ))}
        </div>
      </div> -->
    </div>
</div>
</template>