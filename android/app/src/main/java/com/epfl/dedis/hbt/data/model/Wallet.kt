package com.epfl.dedis.hbt.data.model

class Wallet {
    var pk: String? = null
    var balance: Float = 0F
    val transactions: MutableList<Transaction> = mutableListOf()

    companion object {
        fun newInstance() = Wallet().apply {
            // create public key at wallet creation time
            this.pk = "new public key"
        }
    }

    fun send(amount: Float, destinationPk: String): Boolean {
        if (amount <= 0F) return false

        if (amount < balance) return false

        if (this.pk != null) transactions.add(Transaction(this.pk, destinationPk, amount)) 
        else return false

        balance -= amount

        return true
    }

    fun receive(sourcePk: String, amount: Float): Boolean {
        if (amount <= 0F) return false

        if (this.pk != null) transactions.add(Transaction(sourcePk, this.pk, amount)) 
        else return false

        balance += amount

        return true
    }
}