<template>
    <Cell>
        <div class="big-title">NetworkTest</div>
    </Cell>

    <div style="width: 100%;display: flex;justify-content: center;align-items: center;background-color: white;">
        <div :class="iconClass">
            <div class="icon-inside">
                <h1 :class="networkResultClass">{{ networkResultLable }}</h1>
                <div>network result</div>
            </div>
        </div>
    </div>
    <CellGroup>
        <Cell title="Local Test" :value="testInfo?.local?.label" :value-class="getValueClass(testInfo?.local?.label)">
        </Cell>
        <Cell title="Gateway Test" :value="testInfo?.gateway?.label"
            :value-class="getValueClass(testInfo?.gateway?.label)"></Cell>
        <Cell title="Dns Test" :value="testInfo?.dns?.label" :value-class="getValueClass(testInfo?.dns?.label)"></Cell>
        <Cell title="Internet Test" :value="testInfo?.internet?.label"
            :value-class="getValueClass(testInfo?.internet?.label)"></Cell>
    </CellGroup>
</template>

<script setup>
import * as api from "../api"
import { ref, reactive } from 'vue';
import { CellGroup, Cell, Button, Field, Dialog, Notify, Toast, NoticeBar } from "vant"

let testInfo = reactive({
    local: { label: "testing..." },
    gateway: { label: "testing..." },
    dns: { label: "testing..." },
    internet: { label: "testing..." }
});

let networkResultLable = ref("testing")
let networkResultClass = ref("")
let iconClass = ref("icon-content");

function getValueClass(label) {
    let className = "testing"
    if (label == "pass") {
        className = "net_pass"
    } else if (label == "failed") {
        className = "net_failed"
    }
    return className
}

Toast.loading({
    message: 'newtwork testing...',
    forbidClick: true,
    loadingType: 'spinner',
    duration: 10 * 1000,
    overlay: true
});

api.networkTest().then(d => {
    let allPass = d.every((item) => {
        return item.ok
    })
    if (allPass) {
        networkResultLable.value = "pass"
        networkResultClass.value = "net_pass"
        iconClass.value = "icon-content pass_background"
    } else {
        networkResultLable.value = "failed"
        networkResultClass.value = "net_failed"
        iconClass.value = "icon-content failed_background"
    }

    let obj = {};
    d.forEach((item, index) => {
        if (item.ok) {
            item.label = "pass";
        } else {
            item.label = "failed";
        }

        obj[item.name] = item;
    });
    Object.assign(testInfo, obj)
    Toast.clear()
}).catch(err => {
    Toast.clear()
    console.error(err)
    // Dialog.alert({ message: "network test failed: " + err.getMsg(err) })
})

</script>
<style lang="less" scoped>
.icon-content {
    width: 200px;
    height: 200px;
    background-color: var(--van-cell-value-color);
    border-radius: 50%;
    display: flex;
    justify-content: center;
    align-items: center;

}

.icon-inside {
    width: 180px;
    height: 180px;
    background-color: white;
    border-radius: 50%;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
}

.pass_background {
    background-color: var(--van-success-color);
}

.failed_background {
    background-color: var(--van-danger-color);
}

::v-deep(.net_pass) {
    color: var(--van-success-color);
}

::v-deep(.net_failed) {
    color: var(--van-danger-color);
}
</style>

