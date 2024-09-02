<template>
  <div class="max-w-4xl mx-auto my-10">
    <div class="flex items-center border-b border-teal-500 py-2">
      <input
        class="appearance-none bg-transparent border-none w-full text-gray-700 mr-3 py-1 px-2 leading-tight focus:outline-none"
        type="text"
        v-model="searchText"
        placeholder="Enter search text..."
        @keyup.enter="search"
        aria-label="Search"
      />
      <button
        class="flex-shrink-0 bg-teal-500 hover:bg-teal-700 border-teal-500 hover:border-teal-700 text-sm border-4 text-white py-1 px-2 rounded"
        type="button"
        @click="search"
      >
        Search
      </button>
    </div>
    <div v-if="loading" class="text-center my-2">
      <span class="text-gray-700">Searching...</span>
    </div>
    <div v-if="searchResults.length" class="mt-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="result in searchResults"
        :key="result.id"
        class="bg-white rounded-lg overflow-hidden shadow-lg p-4"
      >
        <h3 class="text-xl text-gray-900 font-bold">{{ result.name }}</h3>
        <p class="text-gray-600 mt-2">Country: {{ result.country }}</p>
        <p class="text-gray-600">City: {{ result.city }}</p>
        <p class="text-gray-600">Address: {{ result.address }}</p>
        <p class="text-gray-600">Website: <a :href="result.websites" class="text-teal-500 hover:text-teal-700" target="_blank">{{ result.websites }}</a></p>
        <p class="text-gray-600">Categories: {{ result.categories}}</p>
      </div>
    </div>
    <p v-else class="text-gray-700 mt-4">No results found.</p>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      searchText: '',
      searchResults: [],
      loading: false
    };
  },
  methods: {
  async search() {
    this.loading = true;
    const response = await fetch(`http://localhost:8080/search?query=${this.searchText}`, {
      method: 'GET', // or 'POST' with appropriate headers and body
    });
    const data = await response.json();
    
    // Check if 'Rows' property exists in the response
    this.searchResults = data.Rows ? data.Rows : []; // Assign empty array if 'Rows' is undefined
    
    this.loading = false;
  }
}
}
</script>

<style>
/* Add this at the top of your <style> section */
body {
  font-family: 'Roboto', sans-serif;
}

h3 {
  font-family: 'Poppins', sans-serif;
}

/* Adjust the existing classes as needed to incorporate the new fonts */
.text-gray-900 {
  font-weight: 600; /* Makes headings slightly bolder */
}

.text-gray-600, .text-gray-700, button, input {
  font-weight: 500; /* Slightly bolder text for better readability */
}
</style>