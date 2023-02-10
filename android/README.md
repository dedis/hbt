# HBT Android Application

This directory contains the Android application of the HBT project.

Its goal is to allows users to extract necessary data from their passport, store them under a username and usethem in transactions.

## Architecture

The application follows the single acivity architecture. In the single activity architecture, the entire application is built around a single activity that manages all the screens within the app.

The main advantages of doing that is the ease of navigation between screens without significant downside.

The application can be split in 3 features :
 - Login
 - Register
 - Wallet

## Login

The login feature allows user to retrieve their wallet from a pair of username/pincode.

## Register

The register feature is designed to create new users in the system. As a result, it needs to scan passports to extract their signature and store it for the created user.

### Extract a passport's signature

Extracting a passport's signature is not an easy task. It is a multi-step process involving both visual scanning and NFC communication.

First, we need to extract the [MRZ data](https://en.wikipedia.org/wiki/Machine-readable_passport#Passport_booklets) needed to produce the passport [BAC key](https://en.wikipedia.org/wiki/Basic_access_control). The needed informations are the passport number, the date of birth and the valid until date.
This can be done automatically with image analysis and pattern matching.

The image analysis is done using google's [MLKit text recognition module](https://developers.google.com/ml-kit/vision/text-recognition). Is it created directly in the [fragment using it](./app/src/main/java/com/epfl/dedis/hbt/ui/register/ScanPassportFragment.kt).

The mrz extraction from text is done [here](./app/src/main/java/com/epfl/dedis/hbt/service/passport/mrz/MRZExtractor.kt).

Then, using the BAC key, the NFC communication with the passport can be established. It is done using the [JMRTD](https://jmrtd.org/) and [Scuba](https://github.com/ugochirico/SCUBA/tree/master/scuba_sc_android).
With that connection, it is possible to extract the signature.

## Wallet

The wallet allows users to create transactions. They are done using QRCodes.

The process is as follow :
1. The receiver create a QRCode containing its id, the timestamp and the amount of HBT to transfer.
2. The sender scans the QRCode and generate a new one adding its own id as the sender. This is the complete transaction.
3. The receiver scans the complete transaction and validates it.

It is implemented usng a state machine which are defined [here](./app/src/main/java/com/epfl/dedis/hbt/data/transaction/TransactionState.kt) and the transitions are managed [here](./app/src/main/java/com/epfl/dedis/hbt/data/transaction/TransactionStateManager.kt)