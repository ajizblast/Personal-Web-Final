function getData() {
  let name = document.getElementById("name").value;
  let email = document.getElementById("email").value;
  let phone = document.getElementById("phone").value;
  let subject = document.getElementById("subject").value;
  let message = document.getElementById("message").value;

  // validation
  if (name == "") {
    return alert("Nama Harus Di Isi");
  } else if (email == "") {
    return alert("Email Harus Di Isi");
  } else if (phone == "") {
    return alert("No Telpon Harus Di Isi");
  } else if (subject == "") {
    return alert("Subject Harus Dipilih");
  } else if (message == "") {
    return alert("Pesan Harus Di Isi");
  }

  const myMail = "ajizblast@gmail.com";
  let a = document.createElement("a");
  a.href = `mailto:${myMail}?subject=${subject}&body= Hallo nama saya ${name}, saya ingin ${message}, bisakah anda menghubungi saya di ${phone}`;
  a.click();

  // console.log(name)
  // console.log(email)
  // console.log(phone)
  // console.log(subject)
  // console.log(message)

  let data = {
    nama: name,
    email: email,
    telp: phone,
    subject: subject,
    pesan: message,
  };

  console.log(data);
}
