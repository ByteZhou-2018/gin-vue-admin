<template>
  <div class="gva-search-box">
    <el-form
      ref="publicForm"
      label-position="right"
      label-width="80px"
      :model="form"
    >
      <el-form-item label="host">
        <el-input v-model="form.host" />
      </el-form-item>
      <el-form-item label="sshPort">
        <el-input v-model="form.port" />
      </el-form-item>
      <el-form-item label="username">
        <el-input v-model="form.username" />
      </el-form-item>
      <el-form-item label="password">
        <el-input type="password" v-model="form.password" />
      </el-form-item>
      <el-form-item>
        <el-button @click="check">环境检测</el-button>
        <el-button @click="stop">停止</el-button>
      </el-form-item>
    </el-form>


    <el-dialog
    v-model="dialogVisible"
    >
      <div class="min-h-96 bg-gray-900">
        <div v-for="(cmd,key) in cmds" :key="key" :class="{
          'text-green-600': cmd.event === 'message',
          'text-red-600': cmd.event === 'done',
          'text-yellow-600': cmd.event === 'info',
        }">
          {{cmd.data}}
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { serverCheck } from '@/api/cloud'
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { fetchEventSource } from '@microsoft/fetch-event-source';


const publicForm = ref(null)
const form = reactive({
  host: '',
  port: '',
  username: '',
  password: '',
})
const ctrl = ref(null)
const cmds = ref([]);
const dialogVisible = ref(false);
const check = async () =>{
  ctrl.value = new AbortController();
  // const res = await serverCheck(form)
  dialogVisible.value = true;
  cmds.value.push({
    event: 'info',
    data: '开始检测'
  });
  fetchEventSource('/api/server/check', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(form),
    signal: ctrl.value.signal,
    onmessage(ev) {
      cmds.value.push(ev);
      console.log(ev)
    }
  });
}

const stop = () =>{
  ctrl.value.abort();
}
</script>
