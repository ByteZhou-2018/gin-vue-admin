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
        <el-button @click="install">安装环境</el-button>
        <el-button @click="stop">停止</el-button>
      </el-form-item>
    </el-form>


    <el-dialog
    v-model="dialogVisible"
    >
      <div class="min-h-96 bg-gray-900">
        <div v-for="(cmd,key) in cmds" :key="key" :class="{
          'text-gray-200': cmd.event === 'message',
          'text-green-600': cmd.event === 'complete',
          'text-red-600': cmd.event === 'fail'
          // 'text-yellow-600': cmd.event === 'info',
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
import { useSSE } from '@/api/sse'


const publicForm = ref(null)
const form = reactive({
  host: '101.43.147.99',
  port: '22',
  username: 'root',
  password: '',
})
const ctrl = ref(null)
const cmds = ref([]);
const dialogVisible = ref(false);

const install =  async () =>{
  const [c,sse] =  useSSE("/server/install",{
    data: form,
    onmessage(ev) {
      switch (ev.event){
        case 'pending':
          if(cmds.value[cmds.value.length - 1]?.event !== 'pending'){
            cmds.value.push({
              event: 'message',
              data: ev.data
            });
          }
          break;
        case 'complete':
          cmds.value.push({
            event: 'complete',
            data: ev.data
          });
          break;
        case 'fail':
          cmds.value.push({
            event: 'fail',
            data: ev.data
          });
          break;
        default:
          cmds.value.push({
            event: 'message',
            data: ev.data
          });
      }
    }
  })
  ctrl.value = c;

  dialogVisible.value = true;
  cmds.value.push({
    event: 'info',
    data: '开始检测'
  });

  await sse();

}
const check = async () =>{
  const [c,sse] =  useSSE("/server/check",{
    data: form,
    onmessage(ev) {
      switch (ev.event){
        case 'pending':
          if(cmds.value[cmds.value.length - 1]?.event !== 'pending'){
            cmds.value.push({
              event: 'message',
              data: ev.data
            });
          }
          break;
        case 'complete':
          cmds.value.push({
            event: 'complete',
            data: ev.data
          });
          break;
        case 'fail':
          cmds.value.push({
            event: 'fail',
            data: ev.data
          });
          break;
        default:
          cmds.value.push({
            event: 'message',
            data: ev.data
          });
      }
    }
  })
  ctrl.value = c;

  dialogVisible.value = true;
  cmds.value.push({
    event: 'info',
    data: '开始检测'
  });

  await sse();

}

const stop = () =>{
  ctrl.value.abort();
}
</script>
