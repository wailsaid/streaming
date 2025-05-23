<script setup lang="ts">

import Button from "@/components/ui/button/Button.vue";
import Calendar from "@/components/ui/calendar/Calendar.vue";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import CardFooter from "@/components/ui/card/CardFooter.vue";
import Checkbox from "@/components/ui/checkbox/Checkbox.vue";
import Input from "@/components/ui/input/Input.vue";
import Popover from "@/components/ui/popover/Popover.vue";
import PopoverContent from "@/components/ui/popover/PopoverContent.vue";
import PopoverTrigger from "@/components/ui/popover/PopoverTrigger.vue";
import Progress from "@/components/ui/progress/Progress.vue";
import RadioGroup from "@/components/ui/radio-group/RadioGroup.vue";
import RadioGroupItem from "@/components/ui/radio-group/RadioGroupItem.vue";
import Select from "@/components/ui/select/Select.vue";
import SelectContent from "@/components/ui/select/SelectContent.vue";
import SelectItem from "@/components/ui/select/SelectItem.vue";
import SelectTrigger from "@/components/ui/select/SelectTrigger.vue";
import SelectValue from "@/components/ui/select/SelectValue.vue";
import Slider from "@/components/ui/slider/Slider.vue";
import Switch from "@/components/ui/switch/Switch.vue";
import Tabs from "@/components/ui/tabs/Tabs.vue";
import TabsContent from "@/components/ui/tabs/TabsContent.vue";
import TabsList from "@/components/ui/tabs/TabsList.vue";
import TabsTrigger from "@/components/ui/tabs/TabsTrigger.vue";
import Textarea from "@/components/ui/textarea/Textarea.vue";
import axios from "axios";
import { CalendarIcon, Crop, Languages, Palette, RotateCcw, Sliders, Sparkles, Upload, Wand2 } from "lucide-vue-next";
import { Label } from "reka-ui";
import { ref } from "vue";

const chunck_size = 5 << 20;

const upload = async (e) => {
  const file = e.target[0].files[0]
  const chunks_count = Math.ceil(file.size / chunck_size)
  for (let i = 0; i < chunks_count; i++) {
    let start = i * chunck_size;
    let end = Math.min(start + chunck_size, file.size);
    
    const chunck = file.slice(start,end)
    const formData = new FormData();

    formData.append("file", chunck);
    formData.append("file_name", file.name);
    formData.append("index", ''+i);
    formData.append("total", chunks_count+'');

    
    try {
      const response = await axios.post("/upload", formData);
      console.log(response.data);
    } catch (error) {
      console.error(error);
      return
    }
  }
};

const isDragging = ref(false);
const file = ref<File | null>(null);
const uploading = ref(false);
const uploadProgress = ref(0);

  const handleDragOver = (e) => {
    e.preventDefault()
    isDragging.value = true
  }
const handleDragLeave = () => {
    isDragging.value = false
  }

  const handleDrop = (e) => {
    e.preventDefault()
    isDragging.value=false

    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      const droppedFile = e.dataTransfer.files[0]
      if (droppedFile.type.startsWith("video/")) {
        file.value = droppedFile
        //generateThumbnail(droppedFile)
      } else {
        alert("Please upload a video file")
      }
    }
  }

    const handleFileChange = (e) => {
    if (e.target.files && e.target.files.length > 0) {
      const selectedFile = e.target.files[0]
      file.value = selectedFile
      console.log(file.value.name);
      
      //generateThumbnail(selectedFile)
    }
  }
</script>
<template>

<div class="container mx-auto py-6">
      <h1 class="text-3xl font-bold mb-6">Upload Video</h1>

      
        <Card v-if="!file" class="w-full max-w-3xl mx-auto">
          <CardHeader>
            <CardTitle>Upload Video</CardTitle>
            <CardDescription>Drag and drop a video file to upload</CardDescription>
          </CardHeader>
          <CardContent>
            <div
              class="border-2 border-dashed rounded-lg p-12 text-center "
                :class="isDragging ? 'border-primary bg-primary/5' : 'border-muted-foreground/20'"
              v-on:dragover="handleDragOver"
              v-on:dragleave="handleDragLeave"
              v-on:drop="handleDrop"
            >
              <Upload class="h-12 w-12 mx-auto text-muted-foreground" />
              <h3 class="mt-4 text-lg font-medium">Drag & drop video file</h3>
              <p class="mt-2 text-sm text-muted-foreground mb-4">MP4, WebM, or MOV. Maximum file size 2GB</p>
              <Input id="video-upload" type="file" accept="video/*" class="hidden" v-on:change="handleFileChange" />
              <label for="video-upload" variant="outline" type="button">
                Select File
              </label>
            </div>
          </CardContent>
        </Card>
      
        <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card class="lg:col-span-2">
            <CardHeader>
              <CardTitle>Video Details</CardTitle>
              <CardDescription>Provide information about your video</CardDescription>
            </CardHeader>
            <CardContent>
              <Tabs defaultValue="details" onValueChange="setCurrentTab">
                <TabsList class="grid w-full grid-cols-5">
                  <TabsTrigger value="details">Details</TabsTrigger>
                  <TabsTrigger value="video">Video</TabsTrigger>
                  <TabsTrigger value="visibility">Visibility</TabsTrigger>
                  <TabsTrigger value="monetization">Monetization</TabsTrigger>
                  <TabsTrigger value="advanced">Advanced</TabsTrigger>
                </TabsList>

                <TabsContent value="details" class="space-y-6 pt-4">
                  <div class="flex justify-between items-center">
                    <h3 class="text-lg font-medium">Basic Information</h3>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick="generateAISuggestions"
                      disabled="isAIProcessing"
                      class="flex items-center gap-1"
                    >
                      <Wand2 class="h-4 w-4" />
                      {isAIProcessing ? "Generating..." : "AI Suggestions"}
                    </Button>
                  </div>

                  {aiSuggestions && (
                    <Card class="bg-muted/50 border-dashed">
                      <CardHeader class="py-3">
                        <CardTitle class="text-sm flex items-center">
                          <Sparkles class="h-4 w-4 mr-2 text-yellow-500" />
                          AI Suggestions
                        </CardTitle>
                      </CardHeader>
                      <CardContent class="space-y-4 py-2">
                        <div class="space-y-1">
                          <div class="flex justify-between">
                            <Label class="text-xs">Title Suggestion</Label>
                            <Button
                              variant="ghost"
                              size="sm"
                              class="h-6 text-xs"
                              onClick="() => applyAISuggestion('title')"
                            >
                              Apply
                            </Button>
                          </div>
                          <p class="text-sm">{aiSuggestions.title}</p>
                        </div>
                        <div class="space-y-1">
                          <div class="flex justify-between">
                            <Label class="text-xs">Description Suggestion</Label>
                            <Button
                              variant="ghost"
                              size="sm"
                              class="h-6 text-xs"
                              onClick="() => applyAISuggestion('description')"
                            >
                              Apply
                            </Button>
                          </div>
                          <p class="text-sm line-clamp-2">{aiSuggestions.description}</p>
                        </div>
                        <div class="space-y-1">
                          <div class="flex justify-between">
                            <Label class="text-xs">Tags Suggestion</Label>
                            <Button
                              variant="ghost"
                              size="sm"
                              class="h-6 text-xs"
                              onClick="() => applyAISuggestion('tags')"
                            >
                              Apply
                            </Button>
                          </div>
                          <p class="text-sm">{aiSuggestions.tags}</p>
                        </div>
                      </CardContent>
                    </Card>
                  )}

                  <div class="space-y-2">
                    <Label htmlFor="title">
                      Title <span class="text-red-500">*</span>
                    </Label>
                    <Input
                      id="title"
                      name="title"
                      placeholder="Add a title that describes your video"
                      value="videoDetails.title"
                      onChange="handleInputChange"
                      required
                    />
                    <p class="text-xs text-muted-foreground">{videoDetails.title.length}/100 characters</p>
                  </div>

                  <div class="space-y-2">
                    <Label htmlFor="description">Description</Label>
                    <Textarea
                      id="description"
                      name="description"
                      placeholder="Tell viewers about your video"
                      class="min-h-32 resize-y"
                      value="videoDetails.description"
                      onChange="handleInputChange"
                    />
                    <p class="text-xs text-muted-foreground">{videoDetails.description.length}/5000 characters</p>
                  </div>

                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="space-y-2">
                      <Label htmlFor="category">Category</Label>
                      <Select
                        value="videoDetails.category"
                        onValueChange="(value) => handleSelectChange('category', value)"
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Select category" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="education">Education</SelectItem>
                          <SelectItem value="entertainment">Entertainment</SelectItem>
                          <SelectItem value="gaming">Gaming</SelectItem>
                          <SelectItem value="music">Music</SelectItem>
                          <SelectItem value="tech">Technology</SelectItem>
                          <SelectItem value="vlog">Vlog</SelectItem>
                          <SelectItem value="howto">How-to & Style</SelectItem>
                          <SelectItem value="science">Science & Technology</SelectItem>
                          <SelectItem value="sports">Sports</SelectItem>
                          <SelectItem value="travel">Travel & Events</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>

                    <div class="space-y-2">
                      <Label htmlFor="language">Language</Label>
                      <Select
                        value="videoDetails.language"
                        onValueChange="(value) => handleSelectChange('language', value)"
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Select language" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="en">English</SelectItem>
                          <SelectItem value="es">Spanish</SelectItem>
                          <SelectItem value="fr">French</SelectItem>
                          <SelectItem value="de">German</SelectItem>
                          <SelectItem value="it">Italian</SelectItem>
                          <SelectItem value="pt">Portuguese</SelectItem>
                          <SelectItem value="ru">Russian</SelectItem>
                          <SelectItem value="ja">Japanese</SelectItem>
                          <SelectItem value="zh">Chinese</SelectItem>
                          <SelectItem value="hi">Hindi</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>

                  <div class="space-y-2">
                    <Label htmlFor="tags">Tags</Label>
                    <Input
                      id="tags"
                      name="tags"
                      placeholder="Add tags (separated by commas)"
                      value="videoDetails.tags"
                      onChange="handleInputChange"
                    />
                    <p class="text-xs text-muted-foreground">
                      Tags help viewers find your video. Use commas to separate tags.
                    </p>
                  </div>

                  <div class="space-y-2">
                    <Label htmlFor="location">Recording Location</Label>
                    <Input
                      id="location"
                      name="location"
                      placeholder="Add a location"
                      value="videoDetails.location"
                      onChange="handleInputChange"
                    />
                  </div>
                </TabsContent>

                <TabsContent value="video" class="space-y-6 pt-4">
                  <div class="space-y-4">
                    <h3 class="text-lg font-medium">Video Elements</h3>

                    <div class="space-y-2">
                      <div class="flex items-center justify-between">
                        <Label htmlFor="endScreen">End Screen</Label>
                        <Switch
                          id="endScreen"
                          checked="videoDetails.endScreen"
                          onCheckedChange="(checked) => handleSwitchChange('endScreen', checked)"
                        />
                      </div>
                      {videoDetails.endScreen && (
                        <Card class="p-4 space-y-3">
                          <p class="text-sm">Choose elements to show in the last 20 seconds of your video</p>
                          {endScreenElements.map((element) => (
                            <div key="element.type" class="flex items-center justify-between">
                              <Label htmlFor="`endscreen-${element.type"`} class="text-sm">
                                {element.type === "subscribe" && "Subscribe button"}
                                {element.type === "related" && "Related videos"}
                                {element.type === "playlist" && "Playlist"}
                              </Label>
                              <Switch
                                id="`endscreen-${element.type"`}
                                checked="element.enabled"
                                onCheckedChange="(checked) => handleEndScreenElementChange(element.type, checked)"
                              />
                            </div>
                          ))}
                        </Card>
                      )}
                    </div>

                    <div class="space-y-2">
                      <div class="flex items-center justify-between">
                        <Label htmlFor="cards">Video Cards</Label>
                        <Switch
                          id="cards"
                          checked="videoDetails.cards"
                          onCheckedChange="(checked) => handleSwitchChange('cards', checked)"
                        />
                      </div>
                      {videoDetails.cards && (
                        <p class="text-sm text-muted-foreground">
                          Cards will be shown throughout your video to promote other content
                        </p>
                      )}
                    </div>

                    <div class="space-y-2">
                      <div class="flex items-center justify-between">
                        <Label>Chapters</Label>
                        <Button variant="outline" size="sm" onClick="handleAddChapter">
                          <Plus class="h-4 w-4 mr-1" /> Add Chapter
                        </Button>
                      </div>
                      <div class="space-y-2">
                        {videoChapters.map((chapter, index) => (
                          <div key="index" class="flex items-center gap-2">
                            <Input
                              placeholder="00:00"
                              value="chapter.time"
                              onChange="(e) => handleChapterChange(index, 'time', e.target.value)"
                              class="w-24"
                            />
                            <Input
                              placeholder="Chapter title"
                              value="chapter.title"
                              onChange="(e) => handleChapterChange(index, 'title', e.target.value)"
                              class="flex-1"
                            />
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick="() => handleRemoveChapter(index)"
                              disabled="videoChapters.length === 1"
                            >
                              <Trash class="h-4 w-4" />
                            </Button>
                          </div>
                        ))}
                      </div>
                      <p class="text-xs text-muted-foreground">
                        Add timestamps and titles to create chapters in your video
                      </p>
                    </div>

                    <div class="space-y-2">
                      <Label>Thumbnail Editor</Label>
                      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div class="space-y-2">
                          <div class="relative aspect-video bg-muted rounded-lg overflow-hidden">
                            {thumbnailPreview ? (
                              <img
                                src=""
                                alt="Video thumbnail"
                                class="w-full h-full object-cover"
                              />
                            ) : (
                              <div class="flex items-center justify-center h-full">
                                <ImageIcon class="h-12 w-12 text-muted-foreground" />
                              </div>
                            )}
                          </div>
                          <div class="flex gap-2">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick="captureCustomThumbnail"
                              disabled="!file"
                              class="flex-1"
                            >
                              <Scissors class="h-4 w-4 mr-1" /> Capture
                            </Button>
                            <Input
                              id="thumbnail-upload"
                              type="file"
                              accept="image/*"
                              class="hidden"
                              onChange="handleThumbnailUpload"
                            />
                            <Button
                              variant="outline"
                              size="sm"
                              onClick="() => document.getElementById('thumbnail-upload')?.click()"
                              class="flex-1"
                            >
                              <Upload class="h-4 w-4 mr-1" /> Upload
                            </Button>
                          </div>
                        </div>
                        <div class="space-y-4">
                          <div class="space-y-2">
                            <Label>Thumbnail Tools</Label>
                            <div class="grid grid-cols-2 gap-2">
                              <Button variant="outline" size="sm" disabled="!thumbnailPreview">
                                <Crop class="h-4 w-4 mr-1" /> Crop
                              </Button>
                              <Button variant="outline" size="sm" disabled="!thumbnailPreview">
                                <RotateCcw class="h-4 w-4 mr-1" /> Rotate
                              </Button>
                              <Button variant="outline" size="sm" disabled="!thumbnailPreview">
                                <Sliders class="h-4 w-4 mr-1" /> Adjust
                              </Button>
                              <Button variant="outline" size="sm" disabled="!thumbnailPreview">
                                <Palette class="h-4 w-4 mr-1" /> Filters
                              </Button>
                            </div>
                          </div>
                          <div class="space-y-2">
                            <Label>Brightness</Label>
                            <Slider :disabled="!thumbnailPreview" :defaultValue="[50]" :max="100" :step="1" />
                          </div>
                          <div class="space-y-2">
                            <Label>Contrast</Label>
                            <Slider :disabled="!thumbnailPreview" :defaultValue="[50]" :max="100" :step="1" />
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  <canvas ref="canvasRef" class="hidden" />
                </TabsContent>

                <TabsContent value="visibility" class="space-y-6 pt-4">
                  <div class="space-y-4">
                    <div class="space-y-2">
                      <Label>Visibility</Label>
                      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <!--   <Card
                          class={`cursor-pointer border-2 ${
                            videoDetails.visibility === "public" ? "border-primary" : "border-muted"
                          }`}
                          onClick="() => handleSelectChange("visibility", "public")"
                        >
                          <CardContent class="flex flex-col items-center justify-center p-4">
                            <Globe class="h-8 w-8 mb-2" />
                            <h3 class="font-medium">Public</h3>
                            <p class="text-xs text-center text-muted-foreground">Everyone can watch</p>
                          </CardContent>
                        </Card>

                        <Card
                          class={`cursor-pointer border-2 ${
                            videoDetails.visibility === "unlisted" ? "border-primary" : "border-muted"
                          }`}
                          onClick="() => handleSelectChange("visibility", "unlisted")"
                        >
                          <CardContent class="flex flex-col items-center justify-center p-4">
                            <Link class="h-8 w-8 mb-2" />
                            <h3 class="font-medium">Unlisted</h3>
                            <p class="text-xs text-center text-muted-foreground">Anyone with the link can watch</p>
                          </CardContent>
                        </Card>

                        <Card
                          class={`cursor-pointer border-2 ${
                            videoDetails.visibility === "private" ? "border-primary" : "border-muted"
                          }`}
                          onClick="() => handleSelectChange("visibility", "private")"
                        >
                          <CardContent class="flex flex-col items-center justify-center p-4">
                            <Lock class="h-8 w-8 mb-2" />
                            <h3 class="font-medium">Private</h3>
                            <p class="text-xs text-center text-muted-foreground">Only you can watch</p>
                          </CardContent>
                        </Card> -->
                      </div>
                    </div>

                    <div class="space-y-2">
                      <Label>Publish Options</Label>
                      <RadioGroup
                        value="videoDetails.publishOption"
                        onValueChange="(value) => handleSelectChange('publishOption', value)"
                      >
                        <div class="flex items-center space-x-2">
                          <RadioGroupItem value="now" id="publish-now" />
                          <Label htmlFor="publish-now">Publish now</Label>
                        </div>
                        <div class="flex items-center space-x-2">
                          <RadioGroupItem value="schedule" id="publish-schedule" />
                          <Label htmlFor="publish-schedule">Schedule for later</Label>
                        </div>
                      </RadioGroup>

                      {videoDetails.publishOption === "schedule" && (
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
                          <div class="space-y-2">
                            <Label>Date</Label>
                            <Popover>
                              <PopoverTrigger asChild>
                                <Button variant="outline" class="w-full justify-start text-left font-normal">
                                  <CalendarIcon class="mr-2 h-4 w-4" />
                                  {scheduledDate ? format(scheduledDate, "PPP") : "Select date"}
                                </Button>
                              </PopoverTrigger>
                              <PopoverContent class="w-auto p-0">
                                <Calendar
                                  mode="single"
                                  selected="scheduledDate"
                                  onSelect="setScheduledDate"
                                  initialFocus
                                />
                              </PopoverContent>
                            </Popover>
                          </div>
                          <div class="space-y-2">
                            <Label>Time</Label>
                            <div class="flex gap-2">
                              <Select defaultValue="12">
                                <SelectTrigger>
                                  <SelectValue placeholder="Hour" />
                                </SelectTrigger>
                                <SelectContent>
                                  {Array.from({ length: 12 }, (_, i) => (
                                    <SelectItem key="i + 1" value="(i + 1).toString()">
                                      {i + 1}
                                    </SelectItem>
                                  ))}
                                </SelectContent>
                              </Select>
                              <Select defaultValue="00">
                                <SelectTrigger>
                                  <SelectValue placeholder="Minute" />
                                </SelectTrigger>
                                <SelectContent>
                                  {Array.from({ length: 60 }, (_, i) => (
                                    <SelectItem key="i" value="i.toString().padStart(2, '0')">
                                      {i.toString().padStart(2, '0')}
                                    </SelectItem>
                                  ))}
                                </SelectContent>
                              </Select>
                              <Select defaultValue="AM">
                                <SelectTrigger>
                                  <SelectValue placeholder="AM/PM" />
                                </SelectTrigger>
                                <SelectContent>
                                  <SelectItem value="AM">AM</SelectItem>
                                  <SelectItem value="PM">PM</SelectItem>
                                </SelectContent>
                              </Select>
                            </div>
                          </div>
                        </div>
                      )}
                    </div>

                    <div class="space-y-4">
                      <h3 class="text-lg font-medium">Audience</h3>

                      <div class="flex items-center justify-between">
                        <div class="space-y-0.5">
                          <Label htmlFor="age-restriction">Age restriction (18+)</Label>
                          <p class="text-xs text-muted-foreground">Mark as adult content</p>
                        </div>
                        <Switch
                          id="age-restriction"
                          checked="videoDetails.ageRestriction"
                          onCheckedChange="(checked) => handleSwitchChange('ageRestriction', checked)"
                        />
                      </div>

                      <div class="space-y-2">
                        <Label>Made for kids?</Label>
                        <RadioGroup defaultValue="no">
                          <div class="flex items-center space-x-2">
                            <RadioGroupItem value="yes" id="kids-yes" />
                            <Label htmlFor="kids-yes">Yes, it's made for kids</Label>
                          </div>
                          <div class="flex items-center space-x-2">
                            <RadioGroupItem value="no" id="kids-no" />
                            <Label htmlFor="kids-no">No, it's not made for kids</Label>
                          </div>
                        </RadioGroup>
                        <p class="text-xs text-muted-foreground">
                          Videos made for kids have limited features and data collection
                        </p>
                      </div>
                    </div>

                    <div class="space-y-4">
                      <h3 class="text-lg font-medium">Comments and Ratings</h3>

                      <div class="flex items-center justify-between">
                        <div class="space-y-0.5">
                          <Label htmlFor="comments">Allow comments</Label>
                          <p class="text-xs text-muted-foreground">Let viewers comment on your video</p>
                        </div>
                        <Switch
                          id="comments"
                          checked="videoDetails.allowComments"
                          onCheckedChange="(checked) => handleSwitchChange('allowComments', checked)"
                        />
                      </div>

                      <div class="flex items-center justify-between">
                        <div class="space-y-0.5">
                          <Label htmlFor="ratings">Allow ratings</Label>
                          <p class="text-xs text-muted-foreground">Let viewers rate your video</p>
                        </div>
                        <Switch
                          id="ratings"
                          checked="videoDetails.allowRatings"
                          onCheckedChange="(checked) => handleSwitchChange('allowRatings', checked)"
                        />
                      </div>
                    </div>
                  </div>
                </TabsContent>

                <TabsContent value="monetization" class="space-y-6 pt-4">
                  <div class="flex items-center justify-between">
                    <div class="space-y-0.5">
                      <Label htmlFor="monetize">Monetize with ads</Label>
                      <p class="text-xs text-muted-foreground">Earn revenue from advertisements</p>
                    </div>
                    <Switch
                      id="monetize"
                      checked="videoDetails.monetize"
                      onCheckedChange="(checked) => handleSwitchChange('monetize', checked)"
                    />
                  </div>

                  {videoDetails.monetize && (
                    <div class="space-y-4">
                      <Card>
                        <CardHeader>
                          <CardTitle>Ad Settings</CardTitle>
                          <CardDescription>Configure how ads appear in your video</CardDescription>
                        </CardHeader>
                        <CardContent class="space-y-4">
                          <div class="flex items-center justify-between">
                            <Label htmlFor="pre-roll">Pre-roll ads</Label>
                            <Switch id="pre-roll" defaultChecked />
                          </div>
                          <div class="flex items-center justify-between">
                            <Label htmlFor="mid-roll">Mid-roll ads</Label>
                            <Switch id="mid-roll" />
                          </div>
                          <div class="flex items-center justify-between">
                            <Label htmlFor="post-roll">Post-roll ads</Label>
                            <Switch id="post-roll" defaultChecked />
                          </div>
                          <div class="flex items-center justify-between">
                            <Label htmlFor="overlay">Overlay ads</Label>
                            <Switch id="overlay" defaultChecked />
                          </div>
                        </CardContent>
                      </Card>

                      <Card>
                        <CardHeader>
                          <CardTitle>Content Suitability</CardTitle>
                          <CardDescription>Ensure your content is suitable for advertisers</CardDescription>
                        </CardHeader>
                        <CardContent class="space-y-4">
                          <div class="space-y-2">
                            <Label>Self Certification</Label>
                            <RadioGroup defaultValue="suitable">
                              <div class="flex items-center space-x-2">
                                <RadioGroupItem value="suitable" id="suitable" />
                                <Label htmlFor="suitable">My content is suitable for all advertisers</Label>
                              </div>
                              <div class="flex items-center space-x-2">
                                <RadioGroupItem value="limited" id="limited" />
                                <Label htmlFor="limited">My content may not be suitable for all advertisers</Label>
                              </div>
                            </RadioGroup>
                          </div>

                          <div class="space-y-2">
                            <Label>Content Declaration</Label>
                            <div class="space-y-2">
                              <div class="flex items-center space-x-2">
                                <Checkbox id="paid-promotion" />
                                <Label htmlFor="paid-promotion" class="text-sm">
                                  Contains paid promotion
                                </Label>
                              </div>
                              <div class="flex items-center space-x-2">
                                <Checkbox id="controversial" />
                                <Label htmlFor="controversial" class="text-sm">
                                  Contains controversial topics
                                </Label>
                              </div>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    </div>
                  )}
                </TabsContent>

                <TabsContent value="advanced" class="space-y-6 pt-4">
                  <div class="space-y-4">
                    <h3 class="text-lg font-medium">Distribution</h3>

                    <div class="space-y-2">
                      <Label>Platforms</Label>
                      <div class="space-y-2">
                        <div class="flex items-center space-x-2">
                          <Checkbox
                            id="platform-youtube"
                            checked="videoDetails.distribution.includes('youtube')"
                            onCheckedChange="(checked) => handleDistributionChange('youtube', checked as boolean)"
                          />
                          <Label htmlFor="platform-youtube">YouTube</Label>
                        </div>
                        <div class="flex items-center space-x-2">
                          <Checkbox
                            id="platform-shorts"
                            checked="videoDetails.distribution.includes('shorts')"
                            onCheckedChange="(checked) => handleDistributionChange('shorts', checked as boolean)"
                          />
                          <Label htmlFor="platform-shorts">YouTube Shorts</Label>
                        </div>
                        <div class="flex items-center space-x-2">
                          <Checkbox
                            id="platform-web"
                            checked="videoDetails.distribution.includes('web')"
                            onCheckedChange="(checked) => handleDistributionChange('web', checked as boolean)"
                          />
                          <Label htmlFor="platform-web">Embed on websites</Label>
                        </div>
                      </div>
                    </div>

                    <div class="space-y-2">
                      <Label>License</Label>
                      <RadioGroup
                        value="videoDetails.license"
                        onValueChange="(value) => handleSelectChange('license', value)"
                      >
                        <div class="flex items-center space-x-2">
                          <RadioGroupItem value="standard" id="license-standard" />
                          <Label htmlFor="license-standard">Standard YouTube License</Label>
                        </div>
                        <div class="flex items-center space-x-2">
                          <RadioGroupItem value="creative" id="license-creative" />
                          <Label htmlFor="license-creative">Creative Commons License</Label>
                        </div>
                      </RadioGroup>
                      <p class="text-xs text-muted-foreground">
                        Creative Commons allows others to reuse and remix your content
                      </p>
                    </div>

                    <div class="space-y-2">
                      <Label>Captions</Label>
                      <div class="flex items-center justify-between">
                        <div class="space-y-0.5">
                          <Label htmlFor="caption">Auto-generate captions</Label>
                          <p class="text-xs text-muted-foreground">Create captions automatically</p>
                        </div>
                        <Switch
                          id="caption"
                          checked="videoDetails.caption"
                          onCheckedChange="(checked) => handleSwitchChange('caption', checked)"
                        />
                      </div>
                      <div class="mt-4">
                        <Button variant="outline" size="sm" disabled="!videoDetails.caption">
                          <Languages class="h-4 w-4 mr-1" /> Upload caption file
                        </Button>
                      </div>
                    </div>

                    <div class="space-y-2">
                      <Label>Recording Date</Label>
                      <Popover>
                        <PopoverTrigger asChild>
                          <Button variant="outline" class="w-full justify-start text-left font-normal">
                            <CalendarIcon class="mr-2 h-4 w-4" />
                            {scheduledDate ? format(scheduledDate, "PPP") : "Select recording date"}
                          </Button>
                        </PopoverTrigger>
                        <PopoverContent class="w-auto p-0">
                          <Calendar mode="single" selected="scheduledDate" onSelect="setScheduledDate" initialFocus />
                        </PopoverContent>
                      </Popover>
                    </div>

                    <div class="space-y-2">
                      <Label>Additional Options</Label>
                      <div class="space-y-2">
                        <div class="flex items-center space-x-2">
                          <Checkbox id="allow-embedding" defaultChecked />
                          <Label htmlFor="allow-embedding" class="text-sm">
                            Allow embedding on external websites
                          </Label>
                        </div>
                        <div class="flex items-center space-x-2">
                          <Checkbox id="notify-subs" defaultChecked />
                          <Label htmlFor="notify-subs" class="text-sm">
                            Notify subscribers
                          </Label>
                        </div>
                        <div class="flex items-center space-x-2">
                          <Checkbox id="product-tagging" />
                          <Label htmlFor="product-tagging" class="text-sm">
                            Enable product tagging
                          </Label>
                        </div>
                      </div>
                    </div>
                  </div>
                </TabsContent>
              </Tabs>
            </CardContent>
            <CardFooter class="flex justify-between">
              <Button variant="outline" onClick="() => router.back()">
                Cancel
              </Button>
              <Button type="submit" onClick="handleSubmit" disabled="uploading || !file || !videoDetails.title">
                {uploading ? "Uploading..." : "Upload Video"}
              </Button>
            </CardFooter>
          </Card>

          <div class="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Video Preview</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="space-y-4">
                  <div class="relative aspect-video bg-muted rounded-lg overflow-hidden">
                    {file && (
                      <video
                        ref="videoRef"
                        src=""
                        class="w-full h-full object-contain"
                        controls
                      />
                    )}
                    <Button
                      variant="ghost"
                      size="icon"
                      class="absolute top-2 right-2 bg-black/50 hover:bg-black/70"
                      onClick="handleRemoveFile"
                    >
                      <X class="h-4 w-4 text-white" />
                    </Button>
                  </div>

                  <div class="space-y-2">
                    <p class="text-sm font-medium">{file?.name}</p>
                    <p class="text-xs text-muted-foreground">
                      {file && `${(file.size / (1024 * 1024)).toFixed(2)} MB`}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Upload Status</CardTitle>
              </CardHeader>
              <CardContent>
                {uploading ? (
                  <div class="space-y-2">
                    <Progress value="uploadProgress" class="h-2" />
                    <div class="flex justify-between text-xs text-muted-foreground">
                      <span>Uploading...</span>
                      <span>{uploadProgress}%</span>
                    </div>
                    <p class="text-xs text-muted-foreground mt-2">
                      {uploadProgress < 100 ? "Uploading video to servers..." : "Processing video, please wait..."}
                    </p>
                  </div>
                ) : (
                  <div class="space-y-4">
                    <div class="flex items-center justify-between">
                      <span class="text-sm">Current tab:</span>
                      <span class="text-sm font-medium capitalize">{currentTab}</span>
                    </div>
                    <div class="flex items-center justify-between">
                      <span class="text-sm">Required fields:</span>
<!--                       <span class="`text-sm font-medium ${videoDetails.title ? "text-green-500" : "text-red-500""`}>
                        {videoDetails.title ? "Complete" : "Title required"}
                      </span> -->
                    </div>
                    <div class="flex items-center justify-between">
                      <span class="text-sm">Visibility:</span>
                      <span class="text-sm font-medium capitalize">{videoDetails.visibility}</span>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="grid grid-cols-2 gap-2">
                  <Button variant="outline" size="sm" class="justify-start">
                    <Info class="h-4 w-4 mr-1" /> Help
                  </Button>
                  <Button variant="outline" size="sm" class="justify-start">
                    <DollarSign class="h-4 w-4 mr-1" /> Monetization
                  </Button>
                  <Button variant="outline" size="sm" class="justify-start">
                    <Music class="h-4 w-4 mr-1" /> Add Music
                  </Button>
                  <Button variant="outline" size="sm" class="justify-start">
                    <Share class="h-4 w-4 mr-1" /> Share
                  </Button>
                  <Button variant="outline" size="sm" class="justify-start" onClick="generateAISuggestions">
                    <Wand2 class="h-4 w-4 mr-1" /> AI Assist
                  </Button>
                  <Button variant="outline" size="sm" class="justify-start">
                    <Tag class="h-4 w-4 mr-1" /> Add Tags
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      
    </div>
</template>

<style scoped></style>