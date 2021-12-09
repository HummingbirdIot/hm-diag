<template>
  <Cell>
    <div class="big-title">Setting</div>
  </Cell>

  <CellGroup title="Github repository proxy">
    <Cell>
      <Field v-model="rp" placeholder="mirror url, support http:// and https://">
        <template #button>
          <Button size="small" type="primary" @click="confirm('gitRepo')">Confirm</Button>
        </template>
      </Field>
    </Cell>
  </CellGroup>

  <CellGroup title="Github release file proxy">
    <Cell>
      <Field v-model="rrp" placeholder="mirror url, support http:// and https://">
        <template #button>
          <Button size="small" type="primary" @click="confirm('gitRelease')">Confirm</Button>
        </template>
      </Field>
    </Cell>
  </CellGroup>
</template>

<script setup>
import { ref } from "vue"
import { CellGroup, Cell, Field, Button, NoticeBar, Toast } from 'vant'

// item: gitRepo or gitRelease
function fetchGitProxy(item) {
  return new Promise((resolve, reject) => {
    fetch("/api/v1/config/proxy?item=" + item)
      .then(r => r.json())
      .then(resolve)
  })
}

const rp = ref("")
const rrp = ref("")

fetchGitProxy("gitRepo").then(rpRe => {
  if (rpRe == null || rpRe.code !== 200) {
    Toast.fail("load git repository proxy error")
  } else {
    rp.value = rpRe.data?.value ? rpRe.data.value : ""
    console.log(rp.value)
  }
})
fetchGitProxy("gitRelease").then(rrpRe => {
  if (rrpRe == null || rrpRe.code !== 200) {
    Toast.fail("load git release file proxy error")
  } else {
    rrp.value = rrpRe.data?.value ? rrpRe.data.value : ""
    console.log(rrp.value)
  }
})

function confirm(item) {
  const v = item == 'gitRepo' ? rp.value : rrp.value
  const u = getUrl(v)
  const type = item == 'gitRepo' ? 'mirror' : 'urlPrefix';
  fetch(`/api/v1/config/proxy?item=${item}`, {
    method: 'POST',
    body: JSON.stringify({ type, value: u })
  }).then(r => r.json())
    .then(r => {
      if (r && r.code == 200) {
        Toast.success("set success")
      } else {
        Toast.fail("set error :" + r?.message)
      }
    })
}

function getUrl(v) {
  const url = new URL(v)
  if (url.protocol !== 'http:' && url.protocol !== 'https:') {
    Toast.fail("wrong url, scheme should be http or https")
    return
  }
  return url.origin + '/'
}

</script>

