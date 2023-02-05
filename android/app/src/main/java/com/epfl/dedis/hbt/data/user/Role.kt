package com.epfl.dedis.hbt.data.user

import androidx.annotation.StringRes
import com.epfl.dedis.hbt.R

enum class Role(@StringRes val roleName: Int) {
    BENEFICIARY(R.string.role_beneficiary),
    MERCHANT(R.string.role_merchant)
}