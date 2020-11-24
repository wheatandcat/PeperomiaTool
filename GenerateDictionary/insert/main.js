const admin = require("firebase-admin");
const serviceAccount = require("../serviceAccount.json");
const dictionary = require("../dictionary.json");

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: "https://peperomia-196da.firebaseapp.com",
});

const firestore = admin.firestore();
const settings = { timestampsInSnapshots: true };
firestore.settings(settings);

console.log(dictionary);

for (let i = 0; i < dictionary.length; i++) {
  const item = {};
  const bigrams = dictionary[i].bigrams;
  for (let j = 0; j < bigrams.length; j++) {
    const key = bigrams[j];
    item[key] = true;
  }
  firestore.collection("version/1/dictionary").add({
    text: dictionary[i].text,
    bigrams: item,
  });
}
