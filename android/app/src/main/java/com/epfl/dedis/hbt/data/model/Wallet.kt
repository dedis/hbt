package com.epfl.dedis.hbt.data.model

class Wallet {
    var pk: String? = null
    var balance: Float = 0F
    private val transactions: MutableList<CompleteTransaction> = mutableListOf()

    companion object {
        fun newInstance() = Wallet().apply {
            // create public key at wallet creation time
            this.pk = "new public key"
        }
    }

    fun send(amount: Float, destinationPk: String): Boolean {
        if (amount <= 0F) return false

        if (amount < balance) return false

        if (this.pk != null) transactions.add(
            CompleteTransaction(
                this.pk!!,
                destinationPk,
                amount,
                0 /*TODO Date time */
            )
        )
        else return false

        balance -= amount

        return true
    }

    fun receive(sourcePk: String, amount: Float): Boolean {
        if (amount <= 0F) return false

        if (this.pk != null) transactions.add(
            CompleteTransaction(
                sourcePk,
                this.pk!!,
                amount,
                0 /*TODO Date time */
            )
        )
        else return false

        balance += amount

        return true
    }
}