<template>
  <div class="d-flex w-100">
    <div class="navbar">
      <h2 class="navbar__title mb-4">All Traces</h2>
      <b-button class="custom-button" v-for="btn in 10"
        >Trace {{ btn }}</b-button
      >
    </div>
    <div class="map">
      <Map />
    </div>
    <div class="navbar"></div>
  </div>
</template>

<script>
export default {
  name: "IndexPage",
};
</script>

<style>
Result Skip Results Iframe EDIT ON html,
body {
  width: 100%;
  height: 100%;
}

body {
  background: linear-gradient(-45deg, #4285f4, #34a854, #fbbc05, #ea4335);
  background-size: 400% 400%;
  animation: gradient 15s ease infinite;
}

@keyframes gradient {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.map {
  width: 70%;
}

.navbar {
  width: 15%;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
}

.navbar__title {
    color: white;
  }

.custom-button {
  min-width: 100px;
  margin-bottom: 12px;
  background-color: rgb(255, 255, 255);
  border: rgb(255, 255, 255);
  color: black;
}

.custom-button:hover {
    background-color: black !important;
    color: white !important;
}

h2 {
  color: green;
}
</style>
