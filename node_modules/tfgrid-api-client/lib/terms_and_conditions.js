async function acceptTermsAndConditions(self, documentLink, documentHash){
    return await self.api.tx.tfgridModule.userAcceptTc(documentLink, documentHash).signAndSend(self.key)
}

module.exports = {
    acceptTermsAndConditions
}