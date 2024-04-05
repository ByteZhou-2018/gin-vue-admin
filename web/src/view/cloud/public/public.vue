<template>
  <div>

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
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-search-box">
      <el-form
          label-position="right"
          label-width="80px"
          :model="publish"
      >
        <el-form-item label="origin">
          <el-input v-model="publish.origin" />
        </el-form-item>
        <el-form-item label="目标环境">
          <el-select v-model="publish.systemType">
            <el-option
                v-for="item in systems"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button @click="zipFunc">本地打包</el-button>
          <el-button @click="deployFunc">上传部署</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-dialog
        v-model="dialogVisible"
        :before-close="Close"
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
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useSSE } from '@/utils/sse'
import { deploy, zip } from '@/api/cloud'


const publicForm = ref(null)
const form = reactive({
  host: '101.43.147.99',
  port: '22',
  username: 'root',
  password: '',
})

const publish = reactive({
  origin: '',
  systemType: ''
})

const systems = [
  {
    value: 'linux',
    label: 'linux'
  },
  {
    value: 'windows',
    label: 'windows'
  },
  {
    value: 'mac',
    label: 'mac'
  }
]

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
              event: 'pending',
              data: "正在执行命令"
            });
          }
          cmds.value[cmds.value.length - 1].data += ".";
          break;
        case 'complete':
          cmds.value.push({
            event: 'complete',
            data: ev.data
          });
          ElMessage.success('执行完成')
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
          ElMessage.success('执行完成')
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
    },
    onclose(){
      console.log("close")
    },
    onerror(err){
      console.log(err)
    }
  })
  ctrl.value = c;

  dialogVisible.value = true;

  await sse();

}
const stop = () =>{
  ctrl.value.abort();
}
const Close = () =>{
  dialogVisible.value = false;
  cmds.value = [];
  stop()
}

const zipFunc = () =>{
  zip(publish)
}
const deployFunc = () =>{
  deploy(form)
}




</script>
