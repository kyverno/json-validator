<template>
  <v-app>
    <v-app-bar title="Kyverno JSON Validator" />
    <v-main class="main">
      <JsonEditorVue id="editor" v-model="value" mode="text" :navigationBar="false" :mainMenuBar="false" />
      <template v-if="validated && !error">
        <v-alert :type="result ? 'success' : 'error'" class="rounded-0" closable v-model="validated">
          Validation {{ result ? 'succeed' : 'failed:' }} <span v-if="!result" v-html="resultMessage" />
        </v-alert>
      </template>
      <template v-if="error">
        <v-alert type="error" class="rounded-0">
          ServerError: error
        </v-alert>
      </template>
      <v-btn color="primary" size="large" class="fab" :loading="loading" @click="validate">Validate</v-btn>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import JsonEditorVue from 'json-editor-vue'
import { ref } from 'vue';
import { resolveAPI } from './api';
import { APIResponse } from './types'

const value = ref<object | string>({
  "foo": "bar"
})

const loading = ref(false)

const validated = ref(false)
const result = ref(false)
const resultMessage = ref()
const error = ref()

const api = resolveAPI()

const validate = () => {
  loading.value = true
  error.value = undefined
  validated.value = false

  let body = value.value
  if (typeof body === 'object') {
    body = JSON.stringify(body)
  }

  fetch(`${api}/validate`, {
    body: body as string,
    method: 'POST',
    mode: 'cors',
    cache: 'no-cache',
    headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json'
    }
  }).then((res) => res.json()).then((content: APIResponse) => {
    validated.value = true
    result.value = content.valid
    resultMessage.value = content.message
  }).catch((err: Error) => {
    console.error(err)
    error.value = err.message
  }).finally(() => {
    loading.value = false
  })
}

</script>

<style>
#editor {
  height: calc(100dvh - 124px);
}

.fab {
  position: fixed !important;
  right: 20px;
  bottom: 100px;
}
</style>