<script setup>
import {PauseCircleOutlined, PlayCircleOutlined} from "@ant-design/icons-vue";
import useAxios from "../support/axios.js";
import {computed, ref} from "vue";
import useMessage from "../support/message.js";

const props = defineProps({
    id: {
        type: String
    },
    enabled: {
        type: Boolean,
        default() {
            return false
        }
    }
})
const emit = defineEmits(["changed"])
const {httpPut} = useAxios()
const {errorMessage} = useMessage()
const enabled = computed({
    get() {
        return props.enabled
    },
    set(v) {

    }
})
const loading = ref(false)

const handleEnableHook = (id) => {
    sendRequest(props.id, "enable")
}
const handleDisableHook = () => {
    sendRequest(props.id, "disable")
}
const sendRequest = (id, type) => {
    loading.value = true
    httpPut(`/webhook/${id}/${type}`).exec().then(({payload}) => {
        emit("changed")
    }).catch(e => {
        errorMessage("Error", e).show()
    }).finally(() => {
        loading.value = false
    })
}
</script>

<template>
    <a-space>
        <a-button type="text" size="large" v-if="enabled" :loading="loading"
                  @click="handleDisableHook">
            <template #icon>
                <PlayCircleOutlined color="primary"/>
            </template>
        </a-button>
        <a-button type="text" danger size="large" v-if="!enabled" :loading="loading"
                  @click="handleEnableHook">
            <template #icon>
                <PauseCircleOutlined/>
            </template>
        </a-button>
    </a-space>
</template>

<style scoped>

</style>
