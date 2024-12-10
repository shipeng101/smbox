<template>
  <n-layout>
    <n-layout-header bordered>
      <div class="header-content">
        <div class="logo">
          <h2>SMBox</h2>
        </div>
        <div class="actions">
          <n-switch v-model:value="isDarkTheme" @update:value="toggleTheme">
            <template #checked>🌙</template>
            <template #unchecked>☀️</template>
          </n-switch>
        </div>
      </div>
    </n-layout-header>
    
    <n-layout-content class="content">
      <n-grid :cols="24" :x-gap="16" :y-gap="16">
        <!-- 系统信息卡片 -->
        <n-grid-item :span="8">
          <system-info-card 
            :cpu="systemInfo.cpu" 
            :memory="systemInfo.memory"
            :uptime="systemInfo.uptime"
          />
        </n-grid-item>

        <!-- 服务状态卡片 -->
        <n-grid-item :span="16">
          <service-status-cards :services="services" />
        </n-grid-item>

        <!-- 节点管理 -->
        <n-grid-item :span="24">
          <node-editor 
            v-if="showNodeEditor"
            :nodes="nodes"
            @update="handleNodeUpdate"
          />
        </n-grid-item>

        <!-- 规则编辑器 -->
        <n-grid-item :span="24">
          <rule-editor 
            v-if="showRuleEditor"
            :rules="rules"
            @update="handleRuleUpdate"
          />
        </n-grid-item>
      </n-grid>
    </n-layout-content>
  </n-layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useTheme } from '../composables/useTheme'
import SystemInfoCard from '../components/SystemInfoCard.vue'
import ServiceStatusCards from '../components/ServiceStatusCards.vue'
import NodeEditor from '../components/NodeEditor.vue'
import RuleEditor from '../components/RuleEditor.vue'

const { isDarkTheme, toggleTheme } = useTheme()

const systemInfo = ref({
  cpu: 0,
  memory: 0,
  uptime: ''
})

const services = ref([])
const nodes = ref([])
const rules = ref([])

const showNodeEditor = ref(false)
const showRuleEditor = ref(false)

// 获取系统信息
const fetchSystemInfo = async () => {
  try {
    const response = await fetch('/api/system/info')
    const data = await response.json()
    systemInfo.value = data
  } catch (error) {
    console.error('Failed to fetch system info:', error)
  }
}

onMounted(() => {
  fetchSystemInfo()
  // 设置定时刷新
  setInterval(fetchSystemInfo, 5000)
})
</script>

<style scoped>
.header-content {
  padding: 0 16px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.content {
  padding: 16px;
}
</style> 