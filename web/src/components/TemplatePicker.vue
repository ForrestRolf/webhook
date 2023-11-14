<script setup>
import {computed, nextTick, ref, watch} from "vue";
import useAxios from "../support/axios.js";
import Gist from "./Gist.vue";
const emit = defineEmits(["selected"])
const props = defineProps({
    lang: {
        type: String,
        default() {
            return "shell"
        }
    }
})
const {httpGet} = useAxios()

const visible = ref(false)
const drawerConfig = computed(() => {
    return {
        width: "50%",
        height: "50%",
        placement: "right",
        closable: true,
    }
})
const title = ref()
const templates = ref([])
const templatesByLang = computed(() => {
    return templates.value.filter(t => {
        return t.language === props.lang
    })
})

const fetchTemplates = () => {
    httpGet('/template').exec().then((data) => {
        templates.value = data.payload
    }).catch(e => {

    })
}

const onClose = () => {
    close()
}

const open = () => {
    visible.value = true
}
const close = () => {
    visible.value = false
}
const onSelected = (template) => {
    emit("selected", template)
    close()
}

watch(visible, () => {
    if(visible.value) {
        fetchTemplates()
    }
})
defineExpose({open, close})
</script>

<template>
    <a-drawer :title="title"
              :placement="drawerConfig.placement"
              :closable="drawerConfig.closable"
              :open="visible"
              :width="drawerConfig.width"
              :height="drawerConfig.height"
              @close="onClose">

        <a-row :gutter="[12, 12]">
            <a-col :span="8" v-for="template in templatesByLang">
                <Gist :template="template" select-mode @gist-click="onSelected"></Gist>
            </a-col>
        </a-row>
    </a-drawer>
</template>

<style scoped>

</style>