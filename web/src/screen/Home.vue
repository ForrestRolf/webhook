<script setup>
import {ref} from "vue";
import {PlusOutlined} from '@ant-design/icons-vue';
import {useRouter} from "vue-router";

const router = useRouter()

const hooks = ref([])

const gotoHook = () => {
    router.push({name: "hooks"})
}
</script>

<template>
    <div class="home">
        <a-row class="new-hook">
            <a-col :span="24">
                <a-button type="primary" @click="gotoHook">
                    <template #icon>
                        <PlusOutlined/>
                    </template>
                    New Hook
                </a-button>
            </a-col>
        </a-row>
        <a-row :gutter="12" class="hooks">
            <a-col :span="4">
                hook name
            </a-col>
            <a-col :span="8">
                triggers
            </a-col>
            <a-col :span="8">
                actions<br/>
            </a-col>
            <a-col :span="2">
                status
            </a-col>
            <a-col :span="2">
                operator
            </a-col>
        </a-row>
    </div>
</template>

<style scoped>

</style>
