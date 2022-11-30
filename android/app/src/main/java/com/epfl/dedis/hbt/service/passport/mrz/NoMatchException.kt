package com.epfl.dedis.hbt.service.passport.mrz

class NoMatchException(
    regex: String,
    text: String
) : Exception("Failed to match regex $regex in $text")
