<script setup>
import {DeleteOutlined, EditOutlined} from "@ant-design/icons-vue"
import {computed, nextTick, onMounted, onUpdated} from "vue";

const emit = defineEmits(["edit", "delete", "read", "gistClick"])
const props = defineProps({
    template: {
        type: Object,
        default() {
            return {}
        }
    },
    selectMode: {
        type: Boolean,
        default() {
            return false
        }
    }
})
const template = computed(() => {
    return props.template
})

const handleEdit = () => {
    emit("edit", template.value)
}
const handleDelete = () => {
    emit("delete", template.value)
}
const handleReadMore = () => {
    emit("read", template.value)
}
const handleClick = () => {
    emit("gistClick", template.value)
}
const highlight = () => {
    nextTick(() => {
        Prism.highlightAll()
    })
}

onUpdated(highlight)
onMounted(highlight)
</script>

<template>
    <a-card class="gist" size="small" @click="handleClick">
        <template #title>
            {{ template.title }}
            <a-tag>{{ template.language }}</a-tag>
        </template>
        <template #extra v-if="!props.selectMode">
            <a-button type="text" @click="handleEdit">
                <template #icon>
                    <EditOutlined color="blue"/>
                </template>
            </a-button>
            <a-popconfirm
                title="Are you sure delete this template?"
                ok-text="Yes"
                cancel-text="No"
                placement="bottom"
                @confirm="handleDelete"
            >
                <a-button type="text">
                    <template #icon>
                        <DeleteOutlined color="danger"/>
                    </template>
                </a-button>
            </a-popconfirm>
        </template>
        <a-card-meta>
            <template #description>{{ template.description }}</template>
        </a-card-meta>

        <a-divider></a-divider>
        <pre>
            <code :class="`language-${template.language}`">{{ template.content }}</code>
        </pre>
        <div class="read-more" @click="handleReadMore" v-if="!props.selectMode">
            <br/>
            Read more
        </div>
    </a-card>
</template>

<style lang="less">
.gist {
    height: 600px;
    overflow: hidden;

    .ant-card-extra {
        display: none;
    }

    &:hover {
        box-shadow: 0 3px 3px rgba(56, 65, 74, .1);

        .ant-card-extra {
            display: inline-block;
        }

        .read-more {
            opacity: 1 !important;
        }
    }

    .ant-card-body {
        position: relative;
        height: 100%;

        .read-more {
            position: absolute;
            left: 0;
            bottom: 30px;
            width: 100%;
            height: 100px;
            text-align: center;
            line-height: 40px;
            cursor: pointer;
            background-image: linear-gradient(to top, rgba(255, 255, 255, 1) 0 70%, rgba(255, 255, 255, 0.5) 80% 100%);
            z-index: 99;
            opacity: 0;
        }
    }
}
</style>
