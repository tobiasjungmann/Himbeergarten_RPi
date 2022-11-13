# import serial
import firebase_admin
from firebase_admin import credentials
from firebase_admin import firestore

# from datetime import datetime

cred = credentials.Certificate("/home/pi/credentials/bewaesserungsanlage-1bb22-firebase-adminsdk-5tuat-e0f60978f7.json")
firebase_admin.initialize_app(cred)
db = firestore.client()

if __name__ == '__main__':
    print("starting")
    db.collection(u'plants').document("testid").set({
        u'graph': dict['graph'] + ';'  + '_' + str(123)
    }, merge=True)
#   ser = serial.Serial('/dev/ttyUSB0', 9600, timeout=1)
#   ser.flush()
#  line = ser.readline().decode('utf-8').rstrip()
# line='100 200 300 400 500'
# print(line)
# humidityValues= line.split()

# now = datetime.now()
# date = now.strftime("%Y")+'.'+now.strftime("%m")+'.'+now.strftime("%d")+':'+now.strftime("%H:%M")
#
#
# users_ref = db.collection(u'plants')
# docs = users_ref.stream()
#
# for doc in docs:
#     #print(f'{doc.id} => {doc.)}')
#     dict=doc.to_dict()
#     computedHumidity=550-(int(humidityValues[dict['usedSensorSlot']])-200)
#     print(dict['usedSensorSlot'])
#     if dict['needsWater']:
#         print("needs watering")
#     db.collection(u'plants').document(doc.id).set({
#         u'humidity': humidityValues[dict['usedSensorSlot']],
#         u'graph': dict['graph'] + ';'+date+'_'+str(computedHumidity)
#     }, merge=True)
