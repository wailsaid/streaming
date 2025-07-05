import { createRouter, createWebHistory } from 'vue-router'
//import App from '../App.vue'
import Upload from '../views/upload.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/home.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/upload.vue'),
    },
    {
      path: '/profile/upload',
      name: 'upload',
      component: () => import('../views/upload.vue'),
    },
{
      path: '/watch/:videoId',
      name: 'watch',
      component: () => import('../views/watch.vue'),
    },

    {
      path: '/profile/history',
      name: 'history',
      component: () => import('../App.vue')
    },
    {
      path: '/profile/watch-later',
      name: 'watch-later',
      component: () => import('../App.vue')
    },
    {
      path: '/profile/liked-videos',
      name: 'liked-videos',
      component: () => import('../App.vue')
    },
    {
      path: '/trending',
      name: 'trending',
      component: () => import('../App.vue')
    },
    {
      path: '/subscriptions',
      name: 'subscriptions',
      component: () => import('../App.vue')
    },
    /*
  {
    path: '/about',
    name: 'about',
    // route level code-splitting
    // this generates a separate chunk (About.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/AboutView.vue'),
  }, */
  ],
})

export default router
