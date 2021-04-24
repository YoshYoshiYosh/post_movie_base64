function setAttributesToVideoElement(videoElement, {mediaToBase64, id, width, height}) {
  videoElement.src      = mediaToBase64
  videoElement.id       = id
  videoElement.width    = width
  videoElement.height   = height
  videoElement.controls = true
}

document.addEventListener('DOMContentLoaded', () => {
  let previewImage = document.getElementById('preview')
  let fileInput = document.getElementById('file')
  fileInput.addEventListener('change', () => {
    const file = fileInput.files[0]
    const fileReader = new FileReader()

    // 読み込みを実行
    fileReader.readAsDataURL(file)

    // 読み込みが成功した場合の処理はこちらに
    fileReader.onload = () => {
      const mediaToBase64 = fileReader.result
      const mediaType = mediaToBase64.slice(0, 15)
      console.log('mediaToBase64：\n' + mediaToBase64)
      console.log('mediaType：\n' + mediaType)

      if (mediaType.match(/video/)) {
        const videoElement = document.createElement("video")
        setAttributesToVideoElement(videoElement, {
          mediaToBase64,
          id: "my-movie",
          width:  200,
          height: 200,
        })
        previewImage.appendChild(videoElement)
      } else if (mediaType.match(/image/)) {
        let image = new Image()
        image.id = 'my-image'
        image.height = 100
        image.src = fileReader.result
        previewImage.appendChild(image)
      } else {
        window.alert('Unknown type media.')
      }
    }

    // 読み込みが失敗した場合の処理はこちらに
    fileReader.onerror = () => {
      console.log('failed')
    }
  })
  const form = document.getElementById('my-form')
  form.addEventListener('submit', async (e) => {
    e.preventDefault()
    const data = {movieBase64: document.getElementById('my-movie').src}
    let response = await fetch('http://localhost:5555/movies/', {
      method: 'POST',
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    }).then(res => res.json())
    console.log('response')
    console.log(response)
  })

  const fetchMovieButton = document.getElementById("fetch-button")
  fetchMovieButton.addEventListener("click", async () => {
    let response = await fetch('http://localhost:5555/movies/1').then(res => res.json())
    console.log('response')
    console.log(response)

    const videoElement = document.createElement("video")
    setAttributesToVideoElement(videoElement, {
      mediaToBase64: response.message,
      id: "my-movie",
      width:  200,
      height: 200,
    })
    previewImage.appendChild(videoElement)
  })
})
