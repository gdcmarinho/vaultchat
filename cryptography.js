async function generateRSAKey() {
  const algorithm = {
      name: "RSA-OAEP",
      modulusLength: 4096,
      publicExponent: new Uint8Array([1, 0, 1]),
      hash: "SHA-256",
  };

  const keyUsages = ["encrypt", "decrypt"];

  const keyPair = await window.crypto.subtle.generateKey(
      algorithm,
      true,
      keyUsages
  );

  const publicKey = await window.crypto.subtle.exportKey("jwk", keyPair.publicKey);
  const privateKey = await window.crypto.subtle.exportKey("jwk", keyPair.privateKey);

  const pemPublicKey = await toPEM(publicKey, "PUBLIC KEY");
  const pemPrivateKey = await toPEM(privateKey, "PRIVATE KEY");

  console.log("Pub key (PEM):");
  console.log(pemPublicKey);

  console.log("Private key (PEM):");
  console.log(pemPrivateKey);

  async function toPEM(key, type) {
      const exportedKey = JSON.stringify(key);
      const pemHeader = `-----BEGIN ${type}-----`;
      const pemFooter = `-----END ${type}-----`;
      const pemBody = window.btoa(exportedKey).match(/.{1,64}/g).join('\n');

      return `${pemHeader}\n${pemBody}\n${pemFooter}`;
  }
}

async function encryptMessage(message, publicKey) {
  const importedPublicKey = await window.crypto.subtle.importKey(
      "jwk",
      publicKey,
      {
          name: "RSA-OAEP",
          hash: "SHA-256",
      },
      false,
      ["encrypt"]
  );

  const encoder = new TextEncoder();
  const encodedMessage = encoder.encode(message);

  const encrypted = await window.crypto.subtle.encrypt(
      { name: "RSA-OAEP" },
      importedPublicKey,
      encodedMessage
  );

  return window.btoa(String.fromCharCode(...new Uint8Array(encrypted)));
}

