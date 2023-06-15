const testimonialData = [
  {
    name: "Martha Cristina Tiahahu",
    quote: "Perang bukan satu-satunya menuju merdeka",
    image: "/public/images/cristina.jpg",
    rating: 2,
  },
  {
    name: "Jendral Sudirman",
    quote: "Mantap sekali ,aku yang pembuat facebook kalah",
    image: "/public/images/jendral-sudirman.jpg",
    rating: 5,
  },
  {
    name: "Bung Karno",
    quote:
      "Masnya ini jago UI/UX, jago slicing juga, database juga mantap. paket lengkap ya mas",
    image: "/public/images/sukarno.jpg",
    rating: 3,
  },
  {
    name: "Bung Tomo",
    quote: "Sepertinya pembuat website ini ingin ku recruit ke perusahaanku",
    image: "/public/images/bung-tomo.jpg",
    rating: 4,
  },
  {
    name: "R.A. Kartini",
    quote: "Pembuat wbsite andalanku dari dulu hingga nanti",
    image: "/public/images/kartini.pg.jpg",
    rating: 5,
  },
  {
    name: "Kapiten Pattimura",
    quote:
      "Selalu ada plus minus, tapi sama mas aji plusnya bisa membuat Indonesia merdeka lagi",
    image: "/public/images/pattimura.png",
    rating: 4,
  },
];

function showTestimonials() {
  let testimonialforHTML = "";

  testimonialData.forEach(function (item) {
    testimonialforHTML += `        
          <div class="testi-card">
           <img src="${item.image}" alt="testi">
           <div class="testi-desc">
               <p class="quotes">"${item.quote}"</p>
               <p class="author">- ${item.name}</p>
               <p class="author">${item.rating} <i class="fa-solid fa-star"></i></p>
           </div>
       </div>`;
  });

  document.getElementById("testimonials").innerHTML = testimonialforHTML;
}

showTestimonials();

function filterTestimonials(rating) {
  let testimonialforHTML = "";

  const testimonialsFiltered = testimonialData.filter(function (item) {
    return item.rating == rating;
  });

  //   console.log(filterTestimonials);

  if (testimonialsFiltered.length === 0) {
    testimonialforHTML += `<h1>Data not found!</h1>`;
  } else {
    testimonialsFiltered.forEach(function (item) {
      testimonialforHTML += `        
              <div class="testi-card">
               <img src="${item.image}" alt="testi">
               <div class="testi-desc">
                   <p class="quotes">"${item.quote}"</p>
                   <p class="author">- ${item.name}</p>
                   <p class="author">${item.rating} <i class="fa-solid fa-star"></i></p>
               </div>
           </div>`;
    });
  }

  document.getElementById("testimonials").innerHTML = testimonialforHTML;
}
